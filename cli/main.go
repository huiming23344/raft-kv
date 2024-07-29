package main

import (
	"fmt"
	kvscli "github.com/huiming23344/kv-raft/client"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	var rootCmd = &cobra.Command{Use: "kvsctl"}
	rootCmd.PersistentFlags().StringP("address", "a", "127.0.0.1:2317", "Server address")
	rootCmd.AddCommand(NewSetCommand(), NewGetCommand(), NewDeleteCommand(), NewMemberCommand())
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func connectServer(cmd *cobra.Command) *kvscli.Client {
	addr, err := cmd.Flags().GetString("address")
	if err != nil {
		log.Fatal(err)
	}
	client, err := kvscli.NewClient(addr)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func NewSetCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "set",
		Short: "Set key to hold the string value",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			rsp, err := connectServer(cmd).Set(args[0], args[1])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(rsp)
		},
	}
	return cc
}

func NewGetCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "get",
		Short: "Get the value of key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rsp, err := connectServer(cmd).Get(args[0])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(rsp)
		},
	}
	return cc
}

func NewDeleteCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "del",
		Short: "Delete the value of key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rsp, err := connectServer(cmd).Del(args[0])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(rsp)
		},
	}
	return cc
}

func NewMemberCommand() *cobra.Command {
	mc := &cobra.Command{
		Use:   "member",
		Short: "Membership related commands",
	}
	mc.AddCommand(NewMemberAddCommand())
	mc.AddCommand(NewMemberRemoveCommand())
	mc.AddCommand(NewMemberListCommand())
	return mc
}

func NewMemberAddCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "add",
		Short: "Adds a member into the cluster",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			rsp, err := connectServer(cmd).Member("add", args[0], args[1])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(rsp)
		},
	}
	return cc
}

func NewMemberRemoveCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "remove",
		Short: "Removes a member from the cluster",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rsp, err := connectServer(cmd).Member("remove", args[0], "")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(rsp)
		},
	}

	return cc
}

func NewMemberListCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "list",
		Short: "Lists all members in the cluster",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			rsp, err := connectServer(cmd).Member("list", "", "")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(rsp)
		},
	}
	return cc
}
