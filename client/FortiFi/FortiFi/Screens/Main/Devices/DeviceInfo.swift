//
//  DeviceInfo.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/22/25.
//

import SwiftUI

struct DeviceInfo: View {
    var device: DevicesResponse
    var body: some View {
        VStack(alignment: .leading) {
            Text("\(device.name)")
                .Title()
                .foregroundStyle(.fortifiForeground)
            
            HStack {
                Text("Date added")
                    .Label()
                    .foregroundStyle(.foregroundMuted)
                Spacer()
                Text(device.date_added)
                    .Label()
                    .foregroundStyle(.fortifiForeground)
            }
            .padding(.horizontal)
            .padding(.vertical, 16)
            .background(.fortifiBackground)
            .cornerRadius(16)
            .shadow(color: Color.black.opacity(0.1), radius: 5, x: 0, y: 2)
            
            VStack(alignment: .leading) {
                HStack {
                    Text("IP address")
                        .Label()
                        .foregroundStyle(.foregroundMuted)
                    Spacer()
                    Text(device.ip_address)
                        .Label()
                        .foregroundStyle(.fortifiForeground)
                }
                Divider()
                HStack {
                    Text("Mac address")
                        .Label()
                        .foregroundStyle(.foregroundMuted)
                    Spacer()
                    Text(device.mac_address)
                        .Label()
                        .foregroundStyle(.fortifiForeground)
                }
            }
            .padding(.horizontal)
            .padding(.vertical, 16)
            .background(.fortifiBackground)
            .cornerRadius(16)
            .shadow(color: Color.black.opacity(0.1), radius: 5, x: 0, y: 2)
            
            HStack {
                Text("Severity")
                    .Label()
                    .foregroundStyle(.foregroundMuted)
                Spacer()
                switch device.incident_count {
                case 0...5:
                    SeverityTag(level: .low)
                case 5...10:
                    SeverityTag(level: .medium)
                default:
                    SeverityTag(level: .high)
                }
            }
            .padding(.horizontal)
            .padding(.vertical, 16)
            .background(.fortifiBackground)
            .cornerRadius(16)
            .shadow(color: Color.black.opacity(0.1), radius: 5, x: 0, y: 2)
            
            Spacer()
        }
        .frame(maxHeight: .infinity)
        .padding()
        .background(.backgroundAlt)
    }
}

enum severity {
    case low, medium, high
    
    var description: String {
        switch self {
        case .low: return "Low"
        case .medium: return "Medium"
        case .high: return "High"
        }
    }
    
    var color: Color {
        switch self {
        case .low: return .fortifiPositive
        case .medium: return .fortifiWarning
        case .high: return .fortifiNegative
        }
    }
    
    var background: Color {
        switch self {
        case .low: return .positiveBackground
        case .medium: return .warningBackground
        case .high: return .negativeBackground
        }
    }
}

struct SeverityTag: View {
    var level: severity
    
    var body: some View {
        Text(level.description)
            .Tag()
            .foregroundStyle(level.color)
            .padding(.horizontal,10)
            .padding(.vertical, 6)
            .background(level.background)
            .cornerRadius(4)
            .overlay(
                RoundedRectangle(cornerRadius: 4)
                    .stroke(.fortifiBorder, lineWidth: 1)
            )
    }
    
}

//#Preview {
//    DeviceInfo()
//}
