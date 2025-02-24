//
//  ChatView.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/23/25.
//

import SwiftUI

struct ChatView: View {
    
    @ObservedObject var viewModel = ChatViewModel()
    
    var body: some View {
        VStack(spacing: 0) {
            Text("AI Chatbot")
                .Header()
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
                                    .background(.textBubble)
                                    .cornerRadius(24)
                            }
                            else {
                                Text("\(message.text)")
                                Spacer()
                            }
                        }
                    }
                }
            }
            TextField("Ask AI a question", text: $viewModel.input)
                .textFieldStyle(CustomTextFieldStyle())
                .background(.fortifiBackground)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(.fortifiBackground)
        .padding()
    }
}

#Preview {
    ChatView()
}
