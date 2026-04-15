package main

import (
	"fmt"
	"os"
)

func main() {
	var rows, seatsAmount int
	fmt.Println("Enter the number of rows:")
	fmt.Scan(&rows)
	fmt.Println("Enter the number of seats in each row:")
	fmt.Scan(&seatsAmount)
	room := createRoom(rows, seatsAmount)

	for {
		menuItem := printMenu()
		switch menuItem {
		case 1:
			room.Print()
		case 2:
			byTicket(&room)
		case 3:
			showStatistics(&room)
		case 0:
			os.Exit(0)
		}
	}
}

func showStatistics(room *Room) {

	statistic := room.GetStatistic()
	fmt.Println()
	fmt.Printf("Number of purchased tickets: %d\n", statistic.Purchased)
	fmt.Printf("Percentage: %.2f%%\n", statistic.Percentage)
	fmt.Printf("Current income: $%d\n", int32(statistic.CurrentIncome))
	fmt.Printf("Total income: $%d\n", int32(statistic.TotalIncome))
	fmt.Println()
}

func byTicket(room *Room) bool {
	fmt.Println("")
	fmt.Println("Enter a row number:")
	var rowNumber, seatNumber int
	fmt.Scan(&rowNumber)
	fmt.Println("Enter a seat number in that row:")
	fmt.Scan(&seatNumber)

	seatExist := room.SeatExists(rowNumber, seatNumber)

	if !seatExist {
		fmt.Println()
		fmt.Println("Wrong input!")
		return byTicket(room)
	}

	booked := room.Book(rowNumber, seatNumber)
	if !booked {
		fmt.Println()
		fmt.Println("That ticket has already been purchased!")
		return byTicket(room)
	}
	seat := room.GetSeat(rowNumber, seatNumber)
	fmt.Printf("\nTicket price: $%d \n", int32(seat.Price))
	return true
}

func printMenu() int {
	var selected int
	fmt.Println("")
	fmt.Println("1. Show the seats")
	fmt.Println("2. Buy a ticket")
	fmt.Println("3. Statistics")
	fmt.Println("0. Exit")
	fmt.Scan(&selected)
	return selected
}

func createRoom(rows, seats int) Room {

	var priceForFirstHalf, priceForSecondHalf float32 = 10, 8
	var firstHalfRows = rows
	totalSeats := rows * seats

	if totalSeats > 60 {
		firstHalfRows = rows / 2
	}

	seatPositions := make(map[SeatPosition]Seat)

	for row := 1; row <= rows; row++ {
		for seat := 1; seat <= seats; seat++ {
			seatPosition := SeatPosition{row, seat}
			price := priceForFirstHalf
			if row > firstHalfRows {
				price = priceForSecondHalf
			}
			seatPositions[seatPosition] = Seat{false, price}
		}
	}

	return Room{
		rows:       rows,
		seatsInRow: seats,
		seats:      seatPositions,
	}
}

type Room struct {
	rows, seatsInRow int
	seats            map[SeatPosition]Seat
}

func (r Room) Print() {
	fmt.Println("")
	fmt.Println("Cinema:")
	fmt.Print("  ")

	for i := 1; i <= r.seatsInRow; i++ {
		fmt.Printf("%d ", i)
	}

	for row := 1; row <= r.rows; row++ {
		fmt.Println()
		fmt.Printf("%d ", row)
		for col := 1; col <= r.seatsInRow; col++ {
			fmt.Printf("%s ", r.GetSeat(row, col).Title())
		}
	}
	fmt.Println("")
}

func (r Room) GetSeat(rowNumber, seatNumber int) Seat {
	return r.seats[SeatPosition{rowNumber, seatNumber}]
}

func (r Room) SeatExists(rowNumber, seatNumber int) bool {
	_, ok := r.seats[SeatPosition{rowNumber, seatNumber}]
	return ok
}

func (r Room) Book(rowNumber, seatNumber int) (success bool) {
	seatPosition := SeatPosition{rowNumber, seatNumber}
	seat := r.GetSeat(rowNumber, seatNumber)

	if seat.Booked {
		return false
	}

	seat.Booked = true
	r.seats[seatPosition] = seat
	return true
}

func (r Room) GetStatistic() Statistic {

	statistic := Statistic{}

	for row := 1; row <= r.rows; row++ {
		for col := 1; col <= r.seatsInRow; col++ {
			seat := r.GetSeat(row, col)
			statistic.TotalIncome += seat.Price
			if seat.Booked {
				statistic.Purchased++
				statistic.CurrentIncome += seat.Price
			}
		}
	}

	total := float32(r.seatsInRow * r.rows)
	purchased := float32(statistic.Purchased)

	statistic.Percentage = (purchased / total) * 100

	return statistic
}

type SeatPosition struct {
	Row, Col int
}

type Seat struct {
	Booked bool
	Price  float32
}

type Statistic struct {
	Purchased     int
	Percentage    float32
	CurrentIncome float32
	TotalIncome   float32
}

func (s Seat) Title() string {
	if s.Booked {
		return "B"
	}

	return "S"
}
