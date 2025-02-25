//
//  ChatView.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/23/25.
//

import SwiftUI

struct ChatView: View {
    
    @ObservedObject var viewModel: ChatViewModel
    
    var body: some View {
        VStack(spacing: 0) {
            Text("AI Chatbot")
                .Header()
                .foregroundStyle(.fortifiForeground)
                .padding(.bottom, 24)
            ScrollView{
                VStack(spacing: 32){
                    ForEach(viewModel.messages, id: \.id) {message in
                        HStack{
                            if message.sender == 0 {
                                Spacer()
                                Text("\(message.text)")
                                    .padding(.horizontal, 16)
                                    .padding(.vertical, 12)
                                    .foregroundStyle(.fortifiForeground)
                                    .background(.textBubble)
                                    .cornerRadius(24).multilineTextAlignment(.trailing)
                                    .Label()
                            }
                            else {
                                Text("\(message.text)")
                                    .foregroundStyle(.fortifiForeground)
                                    .Label()
                                Spacer()
                            }
                        }
                    }
                    if viewModel.isLoading {
                        HStack {
                            ThreeDots()
                            Spacer()
                        }
                    }
                }
            }
            .contentMargins(5, for: .scrollContent)
            .padding(.bottom, 8)
            TextField("Ask AI a question", text: $viewModel.input)
                .textFieldStyle(CustomTextFieldStyle())
                .background(.fortifiBackground)
                .disabled(viewModel.isLoading)
                .onSubmit {
                    Task {
                        switch viewModel.threatSpecified() {
                        case true:
                            if viewModel.messages.last!.id.contains("offerRecommendations") {
                                viewModel.pushUserMessage()
                                await viewModel.handleNeedRecommendationsResponse() }
                            else if viewModel.messages.last!.id.contains("furtherAssistance") {
                                viewModel.pushUserMessage()
                                await viewModel.getMoreAssistance()
                            }
                        case false:
                            await viewModel.getGeneralAssistance()
                        }
                        viewModel.input = ""
                    }
                }
        }
        .padding()
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(.fortifiBackground)
        .onAppear {
            Task {
                switch viewModel.threatSpecified() {
                case true:
                    await viewModel.getThreatAssistance()
                default:
                    break
                }
            }
        }
    }
    
}

// Loading indicator from https://medium.com/@dit.cu/how-to-create-a-loading-animation-with-three-dots-in-swiftui-44ec7ac16cd5

struct ThreeDots: View {
    @State var loading = false
    
    var body: some View {
        HStack(spacing: 20) {
            Circle()
                .fill(.fortifiPrimary)
                .frame(width: 10, height: 10)
                .scaleEffect(loading ? 1.5 : 0.5)
                .animation(.easeInOut(duration: 0.8).repeatForever(autoreverses: true).delay(0.2), value: loading)
            Circle()
                .fill(.fortifiPrimary)
                .frame(width: 10, height: 10)
                .scaleEffect(loading ? 1.5 : 0.5)
                .animation(.easeInOut(duration: 0.8).repeatForever(autoreverses: true).delay(0.2), value: loading)
            Circle()
                .fill(.fortifiPrimary)
                .frame(width: 10, height: 10)
                .scaleEffect(loading ? 1.5 : 0.5)
                .animation(.easeInOut(duration: 0.8).repeatForever(autoreverses: true).delay(0.4), value: loading)
        }
        .onAppear{
            self.loading = true
        }
        .padding(.leading)
    }
}


#Preview {
    ChatView(viewModel: ChatViewModel())
}
