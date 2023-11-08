package main

import (
	"fmt"
)

type User struct {
	Name   string
	Age    int
	Gender string
}

func main() {
	var user User

	fmt.Println("Enter your name:")
	fmt.Scanln(&user.Name)

	fmt.Println("Enter your age:")
	fmt.Scanln(&user.Age)

	fmt.Println("Enter your gender(male/female):")
	fmt.Scanln(&user.Gender)

	fmt.Printf("Name: %s, Age: %d, Gender: %s\n", user.Name, user.Age, user.Gender)

	mainMenu()
}

func mainMenu() {

	speaker := &Speaker{}
	subscriptionContext := NewSubscriptionContext()

	fmt.Println("==================================")
	fmt.Println("Menu:")
	fmt.Println("1. Player")
	fmt.Println("2. Choose a playlist")
	fmt.Println("3. Subscribe to Premium")
	fmt.Println("4. Exit")
	fmt.Println("=================================")

	var choice int
	fmt.Print("Enter your choice: ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		musicPlayerMenu(speaker)
	case 2:
		choosePlaylistMenu()
	case 3:
		subscriptionMenu(subscriptionContext)
	case 4:
		fmt.Println("Exiting the application...")
		return
	default:
		fmt.Println("Invalid choice. Please enter a valid option.")
	}
}

func musicPlayerMenu(speaker SpeakerInterface) {

	fmt.Println("=================================")
	fmt.Println("Music Player Menu:")
	fmt.Println("1. Play")
	fmt.Println("2. Stop")
	fmt.Println("3. Back")
	fmt.Println("=================================")

	var choice int
	fmt.Print("Enter your choice: ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		playCommand := &PlayCommand{Speaker: speaker}
		playCommand.Execute()
		musicPlayerMenu(speaker)
	case 2:
		pauseCommand := &PauseCommand{Speaker: speaker}
		pauseCommand.Execute()
		musicPlayerMenu(speaker)
	case 3:
		mainMenu()
	default:
		fmt.Println("Invalid choice. Please enter a valid option.")
	}
}

func choosePlaylistMenu() {

	playlistFactory := PlaylistFactory{}

	fmt.Println("Choose a Playlist:")
	fmt.Println("1. Pop Hits")
	fmt.Println("2. Rock Classics")
	fmt.Println("3. Jazz Vibes")
	fmt.Println("4. Back")

	var choice int
	fmt.Print("Enter your choice: ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		popHits := playlistFactory.CreatePlaylist("Pop Hits")
		fmt.Printf("You've chosen the playlist: %s\n", popHits.Name)
		mainMenu()
	case 2:
		rockClassics := playlistFactory.CreatePlaylist("Rock Classics")
		fmt.Printf("You've chosen the playlist: %s\n", rockClassics.Name)
		mainMenu()
	case 3:
		jazzVibes := playlistFactory.CreatePlaylist("Jazz Vibes")
		fmt.Printf("You've chosen the playlist: %s\n", jazzVibes.Name)
		mainMenu()
	case 4:
		mainMenu()
	default:
		fmt.Println("Invalid choice. Please enter a valid option.")
	}

}

func subscriptionMenu(subscriptionContext *SubscriptionContext) {

	fmt.Println("=================================")
	fmt.Println("Subscription Menu:")
	fmt.Println("1. Status check")
	fmt.Println("2. Subscribe")
	fmt.Println("3. Unsubscribe")
	fmt.Println("4. Back")
	fmt.Println("=================================")

	var choice int
	fmt.Print("Enter your choice: ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		status := subscriptionContext.CheckStatus()
		fmt.Println(status)
		subscriptionMenu(subscriptionContext)
	case 2:
		status := subscriptionContext
		if !status.IsActive() {
			paymentMenu(subscriptionContext)
		} else {
			status1 := subscriptionContext.Subscribe()
			fmt.Println(status1)
			subscriptionMenu(subscriptionContext)
		}
	case 3:
		status := subscriptionContext.Unsubscribe()
		fmt.Println(status)
		subscriptionMenu(subscriptionContext)
	case 4:
		mainMenu()
	default:
		fmt.Println("Invalid choice. Please enter a valid option.")
	}
}

func paymentMenu(subscriptionContext *SubscriptionContext) {

	fmt.Println("=================================")
	fmt.Println("Payment Menu:")
	fmt.Println("1. Credit Card")
	fmt.Println("2. PayPal")
	fmt.Println("3. Crypto")
	fmt.Println("4. Back")
	fmt.Println("=================================")

	var choice int
	fmt.Print("Enter your choice: ")
	fmt.Scanln(&choice)

	var paymentMethod PaymentMethod

	switch choice {
	case 1:
		paymentMethod = &CreditCardPayment{}
	case 2:
		paymentMethod = &PayPalPayment{}
	case 3:
		paymentMethod = &CryptoPayment{}
	case 4:
		subscriptionMenu(subscriptionContext)
	default:
		fmt.Println("Invalid choice. Please enter a valid option.")
	}
	paymentStatus := paymentMethod.Pay()
	fmt.Println(paymentStatus)

	status := subscriptionContext.Subscribe()
	fmt.Println(status)
	subscriptionMenu(subscriptionContext)
}

//command

type Command interface {
	Execute()
}

type SpeakerInterface interface {
	Play()
	Pause()
}

type Speaker struct {
	IsPlaying bool
}

func (s *Speaker) Play() {
	if !s.IsPlaying {
		fmt.Printf("playing music\n")
		s.IsPlaying = true
	} else {
		fmt.Printf("music is already playing\n")
	}
}

func (s *Speaker) Pause() {
	if s.IsPlaying {
		fmt.Printf("music paused\n")
		s.IsPlaying = false
	} else {
		fmt.Printf("music is already paused\n")
	}
}

type PlayCommand struct {
	Speaker SpeakerInterface
}

func (plc *PlayCommand) Execute() {
	plc.Speaker.Play()
}

type PauseCommand struct {
	Speaker SpeakerInterface
}

func (psc *PauseCommand) Execute() {
	psc.Speaker.Pause()
}

//command end

//factory

type Playlist struct {
	Name string
}

type PlaylistFactory struct{}

func (pf PlaylistFactory) CreatePlaylist(name string) *Playlist {
	return &Playlist{Name: name}
}

//factory end

//state

type State interface {
	Subscribe() string
	Unsubscribe() string
	CheckStatus() string
}

type SubscriptionContext struct {
	state State
}

func NewSubscriptionContext() *SubscriptionContext {
	return &SubscriptionContext{state: &InactiveState{}}
}

func (s *SubscriptionContext) TransitionTo(state State) {
	s.state = state
}

func (s *SubscriptionContext) Subscribe() string {
	status := s.state.Subscribe()
	if s.state.(*InactiveState) != nil {
		s.TransitionTo(&ActiveState{})
	}
	return status
}

func (s *SubscriptionContext) Unsubscribe() string {
	status := s.state.Unsubscribe()
	if s.state.(*ActiveState) != nil {
		s.TransitionTo(&InactiveState{})
	}
	return status
}

func (s *SubscriptionContext) CheckStatus() string {
	return s.state.CheckStatus()
}

func (s *SubscriptionContext) IsActive() bool {
	_, isActive := s.state.(*ActiveState)
	return isActive
}

type ActiveState struct{}

func (as *ActiveState) Subscribe() string {
	return "Already subscribed (Active)"
}

func (as *ActiveState) Unsubscribe() string {
	return "Unsubscribed successfully (Active)"
}

func (as *ActiveState) CheckStatus() string {
	return "You are currently subscribed (Active)"
}

type InactiveState struct{}

func (is *InactiveState) Subscribe() string {
	return "Subscribed successfully (Inactive)"
}

func (is *InactiveState) Unsubscribe() string {
	return "Already unsubscribed (Inactive)"
}

func (is *InactiveState) CheckStatus() string {
	return "You are currently unsubscribed (Inactive)"
}

//state end

//strategy

type PaymentMethod interface {
	Pay() string
}

type CryptoPayment struct{}

func (c *CryptoPayment) Pay() string {
	return "Paid using Crypto"
}

type PayPalPayment struct{}

func (p *PayPalPayment) Pay() string {
	return "Paid using PayPal"
}

type CreditCardPayment struct{}

func (cc *CreditCardPayment) Pay() string {
	return "Paid using Credit Card"
}

//strategy end
