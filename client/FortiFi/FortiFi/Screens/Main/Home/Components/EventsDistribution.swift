//
//  EventsDistribution.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/20/25.
//

import SwiftUI
import Charts

struct EventsDistribution: View {
    
    let data: [(name:String, count: Int, style: Color)] = [
        (name: "normal", count: HomeViewModel.shared.eventCounts.Benign, style: .fortifiPositive),
        (name: "anomaly", count: HomeViewModel.shared.eventCounts.PortScan, style: .fortifiWarning),
        (name: "malicious", count: HomeViewModel.shared.eventCounts.DDoS, style: .fortifiNegative),
    ]
    
    let weekStart = Calendar(identifier: .gregorian).currentWeekBoundary()!.startOfWeek!
        .formatted(date: .numeric, time: .omitted)
    
    var body: some View {
        VStack (spacing: 24){
            HStack {
                Text("Traffic Distribution")
                    .font(.subheadline)
                    .foregroundStyle(.foregroundMuted)
                Spacer()
                Text("\(weekStart) - Present")
                    .font(.subheadline)
                    .foregroundStyle(.foregroundMuted)
            }
            HStack (spacing: 50){
                VStack(alignment: .leading, spacing: 15){
                    VStack (alignment: .leading){
                        Text("^[**\(HomeViewModel.shared.totalEvents)** total event](inflect: true)")
                            .font(.body)
                        Text("This week")
                            .font(.subheadline)
                            .foregroundStyle(.foregroundMuted)
                    }
                    VStack(alignment: .leading, spacing: 16) {
                        HStack {
                            Text("**\(HomeViewModel.shared.eventCounts.Benign)** Benign")
                                .font(.subheadline)
                            Text("\(HomeViewModel.shared.distributions[.benign] ?? 0, specifier: "%.1f")%")
                                .font(.subheadline)
                                .foregroundStyle(.foregroundMuted)
                        }
                        HStack {
                            Text("**\(HomeViewModel.shared.eventCounts.PortScan)** Port Scan")
                                .font(.subheadline)
                            Text("\(HomeViewModel.shared.distributions[.portScan] ?? 0,specifier: "%.1f")%")
                                .font(.subheadline)
                                .foregroundStyle(.foregroundMuted)
                        }
                        HStack {
                            Text("**\(HomeViewModel.shared.eventCounts.DDoS)** DDoS")
                                .font(.subheadline)
                            Text("\(HomeViewModel.shared.distributions[.ddos] ?? 0,specifier: "%.1f")%")
                                .font(.subheadline)
                                .foregroundStyle(.foregroundMuted)
                        }
                    }
                    .padding(.vertical)
                }
                VStack {
                    if HomeViewModel.shared.totalEvents > 0 {
                        Chart {
                            ForEach(data, id: \.name) {type in
                                SectorMark(angle: .value("percent", type.count), angularInset: 1)
                                    .foregroundStyle(type.style)
                                    .cornerRadius(5)
                            }
                        }
                        .frame(height: 150)
                    } else {
                        Chart{
                            SectorMark(angle: .value("percent", 1),
                                       innerRadius: .ratio(0.7),
                                       angularInset: 1)
                            .foregroundStyle(.fortifiPositive)
                        }
                        .frame(height: 150)
                        .chartBackground { chartProxy in
                          GeometryReader { geometry in
                            if let anchor = chartProxy.plotFrame {
                              let frame = geometry[anchor]
                              Text("Nothing to\nReport")
                                .multilineTextAlignment(.center)
                                .font(.caption)
                                .foregroundStyle(.foregroundMuted)
                                .position(x: frame.midX, y: frame.midY)
                            }
                          }
                        }
                    }
                }
            }
        }
        .padding()
        .background(.fortifiBackground)
        .cornerRadius(12)
        .shadow(color: Color.black.opacity(0.1), radius: 5, x: 2, y: 2)
    }
}

#Preview {
    EventsDistribution()
}
