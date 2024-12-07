package main

import (
    "errors"
    "context"
    "os"
    "fmt"
    "bufio"

    log "github.com/obalunenko/logger"
    "github.com/urfave/cli/v2"

    "github.com/wthys/advent-of-code-2024/solver"
    _ "github.com/wthys/advent-of-code-2024/solutions"
)


func onExit(ctx context.Context) cli.AfterFunc {
    return func(c *cli.Context) error {
        return nil
    }
}


func notFound(ctx context.Context) cli.CommandNotFoundFunc {
    return func(c *cli.Context, command string) {
        if _, err := fmt.Fprintf(
            c.App.Writer,
            "Command [%s] not supported.\n Try --help flag to see how to use it\n",
            command,
        ); err != nil {
            log.WithError(ctx, err).Fatal("Failed to print not found message")
        }
    }
}


func cmdRunFlags() []cli.Flag {
    var flags []cli.Flag

    elapsed := cli.BoolFlag{
        Name: "elapsed",
        Aliases: []string{"e"},
        Usage: "Shows elapsed time metric",
        Required: false,
        HasBeenSet: false,
        EnvVars: []string{"ELAPSED"},
    }

    debug := cli.BoolFlag{
        Name: "debug",
        Aliases: []string{"d"},
        Usage: "Shows debug output",
        Required: false,
        HasBeenSet: false,
        EnvVars: []string{"DEBUG"},
    }

    flags = append(flags, &elapsed, &debug)

    return flags
}

func cmdRun(ctx context.Context) cli.ActionFunc {
    return func (c *cli.Context) error {
        if c.Bool("elapsed") {
            ctx = context.WithValue(ctx, "elapsed", true)
        }

        if c.Bool("debug") {
            ctx = context.WithValue(ctx, "debug", true)
        }

        s, err := solver.GetSolver(c.Args().First())
        if err != nil {
            return err
        }

        res, err := solver.Solve(s, bufio.NewReader(os.Stdin), ctx)

        if err != nil {
            return err
        }

        fmt.Println(res)

        return nil
    }
}


func cmdInputFlags() []cli.Flag {
    var flags []cli.Flag

    session := cli.StringFlag{
        Name: "session",
        Aliases: []string{"s"},
        Usage: "AOC Auth session token for getting inputs directly",
        EnvVars: []string{"AOC_SESSION"},
        Required: true,
        HasBeenSet: false,
    }

    flags = append(flags, &session)

    return flags
}

func cmdInput(ctx context.Context) cli.ActionFunc {
    return func (c *cli.Context) error {

        var sess = c.String("session")
        if sess == "" {
            sess = c.String("s")

            if sess == "" {
                return errors.New("no session token provided")
            }
        }

        ctx = context.WithValue(ctx, "session", sess)

        day := c.Args().First()
        if day == "" {
            return errors.New("no puzzle provided")
        }

        input, err := solver.GetInput(ctx, day, sess)

        if err != nil {
            return err
        }

        fmt.Print(string(input[:]))

        return nil
    }
}


func commands(ctx context.Context) []*cli.Command {
    return []*cli.Command{
        {
            Name: "run",
            Usage: `run a specific solution`,
            Action: cmdRun(ctx),
            Flags: cmdRunFlags(),
            SkipFlagParsing: false,
        },
        {
            Name: "input",
            Usage: `get input for a specific day`,
            Action: cmdInput(ctx),
            Flags: cmdInputFlags(),
            SkipFlagParsing: false,
        },
    }
}

var errExit = errors.New("exit is chosen")

func main() {

    ctx := context.Background()

    app := cli.NewApp()
    app.Name = "aoc2024"
    app.Description = "Solutions of puzzles for Advent of Code 2024" +
        " (https://adventofcode.com/2024)"
    app.Usage = `a command line tool for getting solutions for Advent of Code puzzles`
    app.Authors = []*cli.Author{
        {
            Name: "Robbe Thys",
            Email: "robbe.thys@zardof.be",
        },
    }

    app.CommandNotFound = notFound(ctx)
    app.Commands = commands(ctx)
    app.After = onExit(ctx)

    if err := app.Run(os.Args); err != nil {
        if errors.Is(err, errExit) {
            return
        }

        log.WithError(ctx, err).Fatal("Run failed")
    }


}
