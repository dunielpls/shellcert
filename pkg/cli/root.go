package cli

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var (
    rootCmd = &cobra.Command{
        Use:   "shellcert",
        Short: "Toolkit for centralised SSH certificate-based authentication",
    }
    colorsFlag  bool
    verboseFlag bool
    configFile  string
)

func init() {
    // Catch interrupts.
    cancel := make(chan struct{}, 1)
    sigs := make(chan os.Signal, 1)

    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        select {
        case sig := <-sigs:
            fmt.Printf("caught signal %s, exiting\n", sig)
            os.Exit(1)
        case <-cancel:
            fmt.Println("received cancel signal, exiting")
            os.Exit(1)
        }
    }()

    // CLI initialization.
    cobra.EnableCommandSorting = true
    cobra.EnablePrefixMatching = false
    // Disable the default completion` command.
    rootCmd.CompletionOptions.DisableDefaultCmd = true

    // Run `initHook` before each command is executed.
    cobra.OnInitialize(initHook)

    // Root-level Cobra flags.
    rootCmd.PersistentFlags().BoolVarP(&colorsFlag, "no-color", "n", false, "disable colors")
    rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "enable verbose output")
    rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "/etc/shellcert/server.conf", "path to configuration file")
}

func initHook() {
    // Basic configuration.
    viper.SetConfigFile(configFile)
    viper.SetConfigType("json")
    viper.SetEnvPrefix("shellcert")
    viper.AutomaticEnv()

    // Defaults
    viper.SetDefault("server.bind", "0.0.0.0")
    viper.SetDefault("server.port", 8081)
    // TODO: Implement.
    viper.SetDefault("server.trusted_proxies", "")

    // Read config file.
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            fmt.Printf("error parsing configuration file: %v\n", err)
        }
    }
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Printf("fatal error: %v\n", err)
        os.Exit(1)
    }
}

