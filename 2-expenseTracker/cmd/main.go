package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/nycholasmarques/expenseTracker/internal/db"
	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()

	var rootCommand = &cobra.Command{}
	var Description, amount, id, month string

	connStr := "password=root user=postgres dbname=postgres sslmode=disable"
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()
	queries := db.New(dbConn)

	var cmd = &cobra.Command{
		Use:   "criar",
		Short: "Adicione uma despesa",
		Run: func(cmd *cobra.Command, args []string) {

			if Description == "" {
				fmt.Println("Você precisa digitar uma descrição.")
				return
			}
			if amount == "" {
				fmt.Println("Você precisa adicionar um valor.")
				return
			}

			_, err := strconv.Atoi(amount)
			if err != nil {
				fmt.Println("Você precisa passar um número valido")
				return
			}

			description := Description

			expenseParams := db.CreateExpenseParams{
				Description: description,
				Amount:      amount,
				
			}

			expenseCreate, err := queries.CreateExpense(ctx, expenseParams)
			if err != nil {
				fmt.Println("Erro ao criar no banco de dados")
			}

			fmt.Printf("Despesa adicionada com sucesso. (ID: %v)", expenseCreate.ID)

		},
	}

	var cmd2 = &cobra.Command{
		Use:   "listar",
		Short: "Listar despesas",
		Run: func(cmd *cobra.Command, args []string) {
			listExpense, err := queries.ListExpense(ctx)
			if err != nil {
				fmt.Println("Erro ao listar despesas")
			}

			fmt.Println("ID    Description      Amount")
			for _, value := range listExpense {
				fmt.Printf("%v      %s           $%v\n", value.ID, value.Description, value.Amount)
			}
		},
	}

	var cmd3 = &cobra.Command{
		Use:   "total",
		Short: "Total de despesas",
		Run: func(cmd *cobra.Command, args []string) {
			listExpense, err := queries.ListExpense(ctx)
			if err != nil {
				fmt.Println("Erro ao listar despesas")
			}
			var countTracker float64
			for _, value := range listExpense {
				amountConvert, err := strconv.ParseFloat(value.Amount, 64)
				if err != nil {
					fmt.Printf("erro ao converter numero: %v", err)
				}
				countTracker = countTracker + amountConvert
			}
			fmt.Printf("Total de despesas: %v", countTracker)
		},
	}

	var cmd4 = &cobra.Command{
		Use:   "excluir",
		Short: "Excluir despesa",
		Run: func(cmd *cobra.Command, args []string) {
			idConvert, err := strconv.Atoi(id)
			println(id)
			if err != nil {
				fmt.Printf("erro ao converter id: %v", err)
			}
			err = queries.DeleteExpense(ctx, int32(idConvert))
			if err != nil {
				fmt.Printf("erro ao deletar despesa: %v", err)
			}
			fmt.Println("Despesa deletada com sucesso")
		},
	}

	var cmd5 = &cobra.Command{
		Use:   "mes",
		Short: "Total das despesas do mês",
		Run: func(cmd *cobra.Command, args []string) {
			listExpense, err := queries.ListExpense(ctx)
			if err != nil {
				fmt.Println("Erro ao listar despesas")
			}

			monthInt, err := strconv.Atoi(month)
			if err != nil {
				fmt.Printf("erro ao converter mês: %v", err)
			}

			if monthInt < 1 || monthInt > 12 {
				fmt.Println("Número de mês inválido")
				return
			}

			monthIntTime := time.Month(monthInt)

			var countTracker float64
			for _, value := range listExpense {
				expenseMonth := value.CreatedAt.Month()
				if expenseMonth == monthIntTime {
					amountConvert, err := strconv.ParseFloat(value.Amount, 64)
					if err != nil {
						fmt.Printf("erro ao converter numero: %v", err)
					}
					countTracker = countTracker + amountConvert
				}
			}
			if countTracker == 0 {
				fmt.Printf("Não há despesa no mês: %v", monthIntTime)
			} else {
				fmt.Printf("Total de despesas: %v", countTracker)
			}
			
		},
	}

	cmd.Flags().StringVarP(&Description, "descricao", "d", "", "Descrição da despesa")
	cmd.Flags().StringVarP(&amount, "valor", "v", "", "Valor da despesa")
	cmd4.Flags().StringVarP(&id, "id", "i", "", "id da despesa")
	cmd5.Flags().StringVarP(&month, "mes", "m", "", "mês das despesas")

	rootCommand.AddCommand(cmd, cmd2, cmd3, cmd4, cmd5)
	rootCommand.Execute()
}
