//
//  Chat.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/19/25.
//

import SwiftUI

struct Chat: View {
    var viewModel: ChatViewModel
    var body: some View {
        ChatView(viewModel: viewModel)
    }
}

#Preview {
    Chat(viewModel: ChatViewModel())
}
