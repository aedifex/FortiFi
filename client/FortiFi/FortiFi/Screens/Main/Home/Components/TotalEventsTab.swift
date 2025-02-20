//
//  TotalEventsTab.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/19/25.
//

import SwiftUI

struct TotalEventsTab: View {
    @ObservedObject var viewModel = HomeViewModel.shared
    var body: some View {
        VStack {
            HStack {
                VStack (alignment: .leading, spacing: 8) {
                    Text("**\(viewModel.totalEvents)** total events this week")
                            .font(.body)
                    if viewModel.difference == 0 {
                        Text("Same previous week")
                                .font(.subheadline)
                                .foregroundStyle(Color("Foreground-Muted"))
                    } else if viewModel.difference < 0 {
                        Text("\(viewModel.difference) from previous week")
                                .font(.subheadline)
                                .foregroundStyle(Color("Foreground-Negative"))
                    } else {
                        Text("+\(viewModel.difference) from previous week")
                                .font(.subheadline)
                                .foregroundStyle(Color("Foreground-Positive"))
                    }
                }
                .padding()
                Spacer()
            }
        }
        .background(Color(.white))
        .cornerRadius(12)
        .shadow(color: Color.black.opacity(0.1), radius: 5, x: 0, y: 2)
    }
}

#Preview {
    TotalEventsTab()
}
