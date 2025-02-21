//
//  EventTab.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/20/25.
//

import SwiftUI

struct EventTab: View {
    var threat: Event
    var body: some View {
        NavigationLink(destination: SingleEvent(threat: threat)) {
            HStack(spacing: 15){
                VStack (alignment: .leading, spacing: 15){
                    HStack{
                        switch threat.type {
                        case .ddos:
                            DDoSTag()
                        case .portScan:
                            PortScanTag()
                        default:
                            PortScanTag()
                        }
                        Spacer()
                        Text("\(threat.ts)")
                            .font(.subheadline)
                            .foregroundStyle(Color("Foreground-Muted"))
                    }
                    HStack{
                        VStack (alignment: .leading){
                            Text("Source IP")
                                .font(.subheadline)
                                .foregroundStyle(Color("Foreground-Muted"))
                            Text(threat.src)
                        }
                        Spacer()
                        VStack (alignment: .leading){
                            Text("Destination IP")
                                .font(.subheadline)
                                .foregroundStyle(Color("Foreground-Muted"))
                            Text(threat.dst)
                        }
                    }
                }
                Image(systemName: "chevron.right").foregroundColor(Color("Foreground-Muted"))
            }
            .padding(.horizontal,2)
            .padding(.vertical, 8)
        }
    }
}

struct DDoSTag: View {
    var body: some View {
        Text("DDoS")
            .font(.caption)
            .foregroundStyle(Color("Foreground-Negative"))
            .padding(.horizontal,10)
            .padding(.vertical, 6)
            .background(Color("Negative-Accent"))
            .cornerRadius(4)
            .overlay(
                RoundedRectangle(cornerRadius: 4)
                    .stroke(Color("Tag-Border"), lineWidth: 1)
            )
    }
}

struct PortScanTag: View {
    var body: some View {
        Text("PortScan")
            .font(.caption)
            .foregroundStyle(Color("Foreground-Warn"))
            .padding(.horizontal,10)
            .padding(.vertical, 6)
            .background(Color("Background-Warn"))
            .cornerRadius(4)
            .overlay(
                RoundedRectangle(cornerRadius: 4)
                    .stroke(Color("Tag-Border"), lineWidth: 1)
            )
    }
}

#Preview {
    EventTab(threat: Event(id: "123", details: "details hers", ts: "15:04:05", expires: "2006-01-02 15:04:05.999999994", type: .ddos, src: "10.0.0.1", dst: "10.0.0.2"))
}

#Preview {
    PortScanTag()
    DDoSTag()
}
