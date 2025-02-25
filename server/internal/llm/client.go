package llm

import (
	"context"
	"fmt"

	"github.com/aedifex/FortiFi/config"
	"github.com/aedifex/FortiFi/internal/database"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	SYSTEM_MESSAGE = "You are a network security analyst tasked with providing insights into IoT traffic in a home network to a homeowner. "+
				"Based on the machine learning model's prediction, the traffic from IoT device on the network is classified as %v "+
				"where 0 is Normal, 1 is Malicious PartOfAHorizontalPortScan and 2 is Malicious DDoS with a confidence score of %v. "+
				"The model used the following features:\n"+
				"- Feature 1: ts, the timestamp of the connection event.\n"+
				"- Feature 2: uid, a unique identifier for the connection.\n"+
				"- Feature 3: id.orig_h, the source IP address.\n"+
				"- Feature 4: id.orig_p, the source port.\n"+
				"- Feature 5: id.resp_h, the destination IP address.\n"+
				"- Feature 6: id.resp_p, the destination port.\n"+
				"- Feature 7: proto, the network protocol used.\n"+
				"- Feature 8: conn_state, the state of the connection.\n"+
				"- Feature 9: missed_bytes, the number of missed bytes in the connection.\n"+
				"- Feature 10: orig_pkts, the number of packets sent from the source to the destination.\n"+
				"- Feature 11: orig_ip_bytes, the number of IP bytes sent from the source to the destination.\n"+
				"- Feature 12: resp_pkts, the number of packets sent from the destination to the source.\n"+
				"- Feature 13: resp_ip_bytes, the number of IP bytes sent from the destination to the source.\n"+
				"The event details are as follows:"+
				"Source IP: %v\n"+
				"Destination IP: %v\n"+
				"More Details: %v"+
				"Can you provide an explanation or insights into this prediction?"

	GENERAL_ASSISTANCE_MESSAGE = "You are a network security analyst tasked with providing insights into IoT traffic in a home network to a homeowner. The home owner can ask for insights based on specific threats or general information about the network. Your role is to provide a general explanation of the following: %s. Inform the users of the ability to ask for more threat specific information by navigating to the threats list on the FortiFi home page. If the user asks for information not related to the network or security, please inform them that you are not able to provide information on that topic and restate your purpose."
)
type OpenAIClient struct {
	Client *openai.Client
}

func NewOpenAIClient(config *config.Config) *OpenAIClient {
	return &OpenAIClient{
		Client: openai.NewClient(
			option.WithAPIKey(config.OpenAIKey),
		),
	}
}

func (c *OpenAIClient) GetHelpWithThreat(event *database.Event) (string, error) {

	details_message := "Provide brief details on the behavior detected on the network in simple terms without referring to the machine learning model."
	system_message := fmt.Sprintf(SYSTEM_MESSAGE, event.Type, event.Confidence, event.SrcIP, event.DstIP, event.Details)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(system_message),
		openai.UserMessage(details_message),
	}

	response, err := c.generateLLMResponse(messages)
	if err != nil {
		return "", fmt.Errorf("error generating LLM response: %v", err)
	}

	return response, nil

}

func (c *OpenAIClient) GetMoreAssistance(query string, event *database.Event) (string, error) {
	furtherAssistanceMessage := fmt.Sprintf("User mentioned they need assistance with the following: %s. What specific information or help can I provide based on this?", query)

	system_message := fmt.Sprintf(SYSTEM_MESSAGE, event.Type, event.Confidence, event.SrcIP, event.DstIP, event.Details)
	
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(system_message),
		openai.AssistantMessage(furtherAssistanceMessage),
	}

	response, err := c.generateLLMResponse(messages)
	if err != nil {
		return "", fmt.Errorf("error generating LLM response: %v", err)
	}

	return response, nil
	
}

func (c *OpenAIClient) GetRecommendations(event *database.Event) (string, error) {
	recommendationsMessage := "please provide top 3 recommendations to address the threat and explain in simple terms for the home owner."

	system_message := fmt.Sprintf(SYSTEM_MESSAGE, event.Type, event.Confidence, event.SrcIP, event.DstIP, event.Details)
	
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(system_message),
		openai.UserMessage(recommendationsMessage),
	}

	response, err := c.generateLLMResponse(messages)
	if err != nil {
		return "", fmt.Errorf("error generating LLM response: %v", err)
	}
	return response, nil

}

func (c *OpenAIClient) generateLLMResponse(messages []openai.ChatCompletionMessageParamUnion) (string, error) {

	response, err := c.Client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model: openai.F(openai.ChatModelGPT3_5Turbo),
	})

	if err != nil {
		return "", fmt.Errorf("error generating LLM response: %v", err)
	}
	return response.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) GetGeneralAssistance(query string) (string, error) {
	generalAssistanceMessage := "please provide a general explanation of the following: %s"

	system_message := fmt.Sprintf(GENERAL_ASSISTANCE_MESSAGE, query)
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(system_message),
		openai.UserMessage(generalAssistanceMessage),
	}

	response, err := c.generateLLMResponse(messages)
	if err != nil {
		return "", fmt.Errorf("error generating LLM response: %v", err)
	}
	return response, nil
}
