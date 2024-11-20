package events

import (
	"fmt"
	"math"

	"github.com/aws/aws-sdk-go-v2/aws"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/schollz/progressbar/v3"
)

func HydrateDBWithTickets(repo repositories.GameRepositoryInterface, require int, dispatch map[string]int) {
	tokenMap := loadExistingTokens(repo)
	existingCounts := countExistingTickets(repo, dispatch)

	totalExisting := calculateTotalExisting(existingCounts)
	if totalExisting >= require {
		fmt.Printf("%d tickets are already ready\n", totalExisting)
		return
	}

	remaining := require - totalExisting
	ticketsPerPrize := calculateTicketsPerPrize(require, dispatch, existingCounts)
	bar := initializeProgressBar(require, totalExisting)

	generateAndInsertTickets(repo, ticketsPerPrize, remaining, tokenMap, bar)
	fmt.Printf("\n%d tickets are ready\n", require)
}

func loadExistingTokens(repo repositories.GameRepositoryInterface) map[string]bool {
	tokenMap := make(map[string]bool)
	existingTokens, err := repo.ReadTickets(&transfert.Ticket{})
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch existing tokens: %v", err))
	}

	for _, existingToken := range existingTokens {
		tokenMap[existingToken.Token.String()] = true
	}
	fmt.Printf("Loaded %d existing tokens from the database\n", len(tokenMap))
	return tokenMap
}

func countExistingTickets(repo repositories.GameRepositoryInterface, dispatch map[string]int) map[string]int {
	existingCounts := make(map[string]int)
	for prize := range dispatch {
		count, err := repo.CountTicket(&transfert.Ticket{
			Prize: aws.String(prize),
		})
		if err != nil {
			panic(fmt.Sprintf("Failed to count tickets for %s: %v", prize, err))
		}
		existingCounts[prize] = count
	}
	return existingCounts
}

func calculateTotalExisting(existingCounts map[string]int) int {
	total := 0
	for _, count := range existingCounts {
		total += count
	}
	return total
}

func calculateTicketsPerPrize(require int, dispatch map[string]int, existingCounts map[string]int) map[string]int {
	ticketsPerPrize := make(map[string]int)
	for prize, percent := range dispatch {
		expected := int(math.Round(float64(require) * float64(percent) / 100.0))
		ticketsPerPrize[prize] = expected - existingCounts[prize]
		fmt.Println(ticketsPerPrize[prize], "tickets for", prize, existingCounts[prize], "already exist")
	}
	return ticketsPerPrize
}

func initializeProgressBar(require, totalExisting int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(require,
		progressbar.OptionSetDescription("Inserting tickets..."),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionEnableColorCodes(true),
	)
	bar.Add(totalExisting)
	return bar
}

func generateAndInsertTickets(repo repositories.GameRepositoryInterface, ticketsPerPrize map[string]int, remaining int, tokenMap map[string]bool, bar *progressbar.ProgressBar) {
	modulo := 1000
	generateUniqueToken := func() *string {
		for {
			token := token.Generate(12).PointerString()
			if !tokenMap[*token] {
				tokenMap[*token] = true
				return token
			}
		}
	}

	for prize, numTickets := range ticketsPerPrize {
		if numTickets <= 0 {
			continue
		}
		tickets := []*transfert.Ticket{}
		for i := 0; i < numTickets; i++ {
			tickets = append(tickets, &transfert.Ticket{
				Prize: aws.String(prize),
				Token: generateUniqueToken(),
			})

			if len(tickets) >= modulo || i == numTickets-1 {
				if err := repo.CreateTickets(tickets); err != nil {
					panic(fmt.Sprintf("Failed to insert tickets for %s: %v", prize, err))
				}
				bar.Add(len(tickets))
				tickets = []*transfert.Ticket{}
			}
		}
	}
}
