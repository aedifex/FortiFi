//
//  EventsDistribution.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/20/25.
//

import SwiftUI
import Charts

struct EventsDistribution: View {
    
    let data = [
        (name: "normal", count: HomeViewModel.shared.eventCounts.Normal, style: Color("Foreground-Positive")),
        (name: "anomaly", count: HomeViewModel.shared.eventCounts.Anomalous, style: .yellow),
        (name: "malicious", count: HomeViewModel.shared.eventCounts.Malicious, style: Color("Foreground-Negative")),
    ]
    
    var body: some View {
        HStack (spacing: 60){
            VStack(alignment: .leading, spacing: 15){
                VStack (alignment: .leading){
                    Text("^[**\(HomeViewModel.shared.totalEvents)** total event](inflect: true)")
                        .font(.body)
                    Text("In the past week")
                        .font(.subheadline)
                        .foregroundStyle(Color("Foreground-Muted"))
                }
                VStack(alignment: .leading, spacing: 16) {
                    HStack {
                        Text("**\(HomeViewModel.shared.eventCounts.Normal)** Normal")
                            .font(.subheadline)
                        Text("\(HomeViewModel.shared.distributions[0], specifier: "%.1f")%")
                            .font(.subheadline)
                            .foregroundStyle(Color("Foreground-Muted"))
                    }
                    HStack {
                        Text("**\(HomeViewModel.shared.eventCounts.Anomalous)** Anomaly")
                            .font(.subheadline)
                        Text("\(HomeViewModel.shared.distributions[1],specifier: "%.1f")%")
                            .font(.subheadline)
                            .foregroundStyle(Color("Foreground-Muted"))
                    }
                    HStack {
                        Text("**\(HomeViewModel.shared.eventCounts.Normal)** Malicious")
                            .font(.subheadline)
                        Text("\(HomeViewModel.shared.distributions[2],specifier: "%.1f")%")
                            .font(.subheadline)
                            .foregroundStyle(Color("Foreground-Muted"))
                    }
                }
                .padding(.vertical)
            }
            VStack {
                Chart {
                    ForEach(data, id: \.name) {type in
                        SectorMark(angle: .value("percent", type.count), angularInset: 1)
                            .foregroundStyle(type.style)
                            .cornerRadius(5)
                    }
                }
                .frame(height: 150)
            }
        }
        .padding()
        .background(Color(.white))
        .cornerRadius(12)
        .shadow(color: Color.black.opacity(0.1), radius: 5, x: 2, y: 2)
    }
}

#Preview {
    EventsDistribution()
}
