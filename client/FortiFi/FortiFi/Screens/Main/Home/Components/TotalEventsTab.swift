//
//  TotalEventsTab.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/19/25.
//

import SwiftUI

struct TotalEventsTab: View {
    @ObservedObject var viewModel = HomeViewModel.shared
    let weekStart = Calendar(identifier: .gregorian).currentWeekBoundary()!.startOfWeek!
        .formatted(date: .numeric, time: .omitted)

    var body: some View {
        VStack(alignment: .leading, spacing: 24) {
            HStack {
                Text("Traffic Volume")
                    .font(.subheadline)
                    .foregroundStyle(Color("Foreground-Muted"))
                Spacer()
                Text("\(weekStart) - Present")
                    .font(.subheadline)
                    .foregroundStyle(Color("Foreground-Muted"))
            }
            HStack {
                VStack (alignment: .leading, spacing: 8) {
                    Text("^[**\(viewModel.totalEvents)** total event](inflect: true) this week")
                            .font(.body)
                    if viewModel.difference == 0 {
                        Text("Same as previous week")
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
                Spacer()
            }
        }
        .padding()
        .background(Color(.white))
        .cornerRadius(12)
        .shadow(color: Color.black.opacity(0.1), radius: 5, x: 2, y: 2)
    }
}

#Preview {
    TotalEventsTab()
}
