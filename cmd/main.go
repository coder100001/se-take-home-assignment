package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/feedme/order-controller/internal/controller"
	"github.com/feedme/order-controller/internal/order"
)

func timestamp() string {
	return time.Now().Format("15:04:05")
}

func printHelp() {
	fmt.Println("McDonald's Order Controller")
	fmt.Println("Commands:")
	fmt.Println("  normal <desc>  - Create a normal order")
	fmt.Println("  vip <desc>     - Create a VIP order")
	fmt.Println("  +bot           - Add a cooking bot")
	fmt.Println("  -bot           - Remove the newest bot")
	fmt.Println("  status         - Show current status")
	fmt.Println("  help           - Show this help")
	fmt.Println("  quit           - Exit")
}

func main() {
	ctrl := controller.NewController()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		ctrl.Shutdown()
		fmt.Printf("\n[%s] Goodbye!\n", timestamp())
		os.Exit(0)
	}()

	printHelp()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		cmd := parts[0]
		arg := ""
		if len(parts) > 1 {
			arg = parts[1]
		}

		switch cmd {
		case "normal":
			o := ctrl.CreateOrder(order.Normal, arg)
			fmt.Printf("[%s] %s created\n", timestamp(), o.String())
		case "vip":
			o := ctrl.CreateOrder(order.VIP, arg)
			fmt.Printf("[%s] %s created\n", timestamp(), o.String())
		case "+bot":
			b := ctrl.AddBot()
			fmt.Printf("[%s] Bot#%d added\n", timestamp(), b.ID)
		case "-bot":
			b := ctrl.RemoveBot()
			if b == nil {
				fmt.Printf("[%s] No bots to remove\n", timestamp())
			} else {
				fmt.Printf("[%s] Bot#%d removed\n", timestamp(), b.ID)
			}
		case "status":
			pending, processing, complete, bots := ctrl.GetStatus()
			fmt.Printf("[%s] === Status ===\n", timestamp())

			sort.Slice(pending, func(i, j int) bool {
				return pending[i].Type > pending[j].Type
			})

			var pendingStrs []string
			for _, o := range pending {
				pendingStrs = append(pendingStrs, o.String())
			}
			fmt.Printf("[%s] PENDING: %s\n", timestamp(), strings.Join(pendingStrs, ", "))

			orderBotMap := make(map[int]int)
			for _, b := range bots {
				if b.Order != nil {
					orderBotMap[b.Order.ID] = b.ID
				}
			}

			var procStrs []string
			for _, o := range processing {
				if botID, ok := orderBotMap[o.ID]; ok {
					procStrs = append(procStrs, fmt.Sprintf("%s [Bot#%d]", o.String(), botID))
				} else {
					procStrs = append(procStrs, o.String())
				}
			}
			fmt.Printf("[%s] PROCESSING: %s\n", timestamp(), strings.Join(procStrs, ", "))

			var compStrs []string
			for _, o := range complete {
				compStrs = append(compStrs, o.String())
			}
			fmt.Printf("[%s] COMPLETE: %s\n", timestamp(), strings.Join(compStrs, ", "))

			var botStrs []string
			for _, b := range bots {
				if b.Order != nil {
					botStrs = append(botStrs, fmt.Sprintf("Bot#%d (PROCESSING: Order#%d)", b.ID, b.Order.ID))
				} else {
					botStrs = append(botStrs, fmt.Sprintf("Bot#%d (IDLE)", b.ID))
				}
			}
			fmt.Printf("[%s] Bots: %s\n", timestamp(), strings.Join(botStrs, ", "))

		case "help":
			printHelp()
		case "quit":
			ctrl.Shutdown()
			fmt.Printf("[%s] Goodbye!\n", timestamp())
			return
		default:
			fmt.Printf("[%s] Unknown command. Type 'help' for available commands.\n", timestamp())
		}
	}

	ctrl.Shutdown()
	fmt.Printf("[%s] Goodbye!\n", timestamp())
}
