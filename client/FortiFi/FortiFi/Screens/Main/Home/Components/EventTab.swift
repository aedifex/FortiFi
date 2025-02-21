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
                            .Label()
                            .foregroundStyle(.foregroundMuted)
                    }
                    HStack{
                        VStack (alignment: .leading){
                            Text("Source IP")
                                .Label()
                                .foregroundStyle(.foregroundMuted)
                            Text(threat.src)
                        }
                        Spacer()
                        VStack (alignment: .leading){
                            Text("Destination IP")
                                .Label()
                                .foregroundStyle(.foregroundMuted)
                            Text(threat.dst)
                        }
                    }
                }
                Image(systemName: "chevron.right").foregroundColor(.foregroundMuted)
            }
            .padding(.horizontal,2)
            .padding(.vertical, 8)
        }
    }
}

struct DDoSTag: View {
    var body: some View {
        Text("DDoS")
            .Tag()
            .foregroundStyle(.fortifiNegative)
            .padding(.horizontal,10)
            .padding(.vertical, 6)
            .background(.negativeBackground)
            .cornerRadius(4)
            .overlay(
                RoundedRectangle(cornerRadius: 4)
                    .stroke(.fortifiBorder, lineWidth: 1)
            )
    }
}

struct PortScanTag: View {
    var body: some View {
        Text("PortScan")
            .Tag()
            .foregroundStyle(.fortifiWarning)
            .padding(.horizontal,10)
            .padding(.vertical, 6)
            .background(.warningBackground)
            .cornerRadius(4)
            .overlay(
                RoundedRectangle(cornerRadius: 4)
                    .stroke(.fortifiBorder, lineWidth: 1)
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
