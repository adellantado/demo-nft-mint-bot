package main

import (
	"context"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	END string = "end"
)

type ConversationManager struct {
	registeredConvos map[string]*Conversation
	activeConvos     map[int64]*Conversation
}

func NewConversationManager() *ConversationManager {
	return &ConversationManager{
		registeredConvos: make(map[string]*Conversation),
		activeConvos:     make(map[int64]*Conversation),
	}
}

func (cm *ConversationManager) AddConvo(convoName string, convo *Conversation) {
	cm.registeredConvos[convoName] = convo
}

func (cm *ConversationManager) AddConvoHandlers(convos map[string][]func(context.Context, *bot.Bot, *models.Update) string) {
	for convoName, hadlers := range convos {
		newConvo := NewConversation("0")
		for i, hadler := range hadlers {
			newConvo.AddHandler(strconv.Itoa(i), hadler)
		}
		cm.AddConvo(convoName, newConvo)
	}
}

func (cm *ConversationManager) InitConvo(chatID int64, convoName string) {
	activeConvo, exists := cm.registeredConvos[convoName]
	if exists {
		cm.activeConvos[chatID] = activeConvo
	}
}

func (cm *ConversationManager) Handle(ctx context.Context, b *bot.Bot, update *models.Update) bool {
	chatID := update.Message.Chat.ID
	activeConvo, exists := cm.activeConvos[chatID]
	if exists {
		state := activeConvo.HandleUpdate(ctx, b, update)
		if state == END {
			cm.activeConvos[chatID] = nil
		}
		return true
	}
	return false
}

// Conversation defines a state-based conversation handler
type Conversation struct {
	stateHandlers map[string]func(context.Context, *bot.Bot, *models.Update) string // State functions returning the next state
	userStates    map[int64]string                                                  // Keeps track of each user's current state
	defaultState  string                                                            // Default state to fall back to
}

// NewConversation creates a new Conversation instance
func NewConversation(defaultState string) *Conversation {
	return &Conversation{
		stateHandlers: make(map[string]func(context.Context, *bot.Bot, *models.Update) string),
		userStates:    make(map[int64]string),
		defaultState:  defaultState,
	}
}

// AddHandler adds a handler function for a specific state
func (c *Conversation) AddHandler(state string, handler func(context.Context, *bot.Bot, *models.Update) string) {
	c.stateHandlers[state] = handler
}

// HandleUpdate processes an update and routes it to the correct state handler
func (c *Conversation) HandleUpdate(ctx context.Context, b *bot.Bot, update *models.Update) string {
	// Get the user's current state; if none, use the default state
	chatID := update.Message.Chat.ID
	currentState, exists := c.userStates[chatID]
	if !exists {
		currentState = c.defaultState
	}

	// Get the handler for the current state
	handler, handlerExists := c.stateHandlers[currentState]
	if !handlerExists {
		// If no handler exists for the state, reset to the default state
		c.userStates[chatID] = c.defaultState
		return END
	}

	// Call the handler and get the next state
	nextState := handler(ctx, b, update)

	if nextState == END {
		c.userStates[chatID] = c.defaultState
	} else {
		// Update the user's state
		c.userStates[chatID] = nextState
	}
	return nextState
}

// ResetState resets the conversation state for a specific user
func (c *Conversation) ResetState(chatID int64) {
	c.userStates[chatID] = c.defaultState
}
