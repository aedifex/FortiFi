//
//  SingleEvent.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/20/25.
//

import SwiftUI

struct SingleEvent: View {
    var threat: Event
    @State var offset: CGSize = CGSize(width: 0, height: 600)
    var body: some View {
        VStack (alignment: .leading, spacing: 24){
            Text("Event Details")
                .font(.title)
                .foregroundStyle(Color("Foreground"))
            
            VStack {
                HStack {
                    Text("Source Ip")
                        .font(.subheadline)
                        .foregroundStyle(Color("Foreground-Muted"))
                    Spacer()
                    Text(threat.src)
                        .font(.subheadline)
                        .foregroundStyle(Color("Foreground"))
                }
                .padding(.vertical,4)
                Divider()
                HStack {
                    Text("Destination Ip")
                        .font(.subheadline)
                        .foregroundStyle(Color("Foreground-Muted"))
                    Spacer()
                    Text(threat.dst)
                        .font(.subheadline)
                        .foregroundStyle(Color("Foreground"))
                }
                .padding(.vertical,4)
            }
            .padding()
            .background(.white)
            .cornerRadius(16)
            .shadow(color: Color.black.opacity(0.1), radius: 5, x: 0, y: 2)
            
            HStack {
                Text("Time of Incident")
                    .font(.subheadline)
                    .foregroundStyle(Color("Foreground-Muted"))
                Spacer()
                Text(threat.ts)
                    .font(.subheadline)
                    .foregroundStyle(Color("Foreground"))
            }
            .padding(.horizontal)
            .padding(.vertical, 16)
            .background(.white)
            .cornerRadius(16)
            .shadow(color: Color.black.opacity(0.1), radius: 5, x: 0, y: 2)
            
            HStack {
                Text("Attack Type")
                    .font(.subheadline)
                    .foregroundStyle(Color("Foreground-Muted"))
                Spacer()
                if threat.type == .portScan {
                    PortScanTag()
                } else {
                    DDoSTag()
                }
            }
            .padding(.horizontal)
            .padding(.vertical, 16)
            .background(.white)
            .cornerRadius(16)
            .shadow(color: Color.black.opacity(0.1), radius: 5, x: 0, y: 2)
            
            VStack(alignment: .leading) {
                Text("Details")
                    .font(.subheadline)
                    .foregroundStyle(Color("Foreground-Muted"))
                Divider()
                Text(threat.details)
                    .font(.subheadline)
                    .foregroundStyle(Color("Foreground"))
                    .padding(.vertical)
            }
            .padding(.horizontal)
            .padding(.vertical, 16)
            .background(.white)
            .cornerRadius(16)
            .shadow(color: Color.black.opacity(0.1), radius: 5, x: 0, y: 2)
            
            HStack {
                Spacer()
                NavigationLink(destination: Chat()) {
                    Text("Ask AI for Assistance")
                        .font(.headline)
                        .padding()
                        .foregroundStyle(.white)
                        .background(Color("Fortifi-Primary"))
                }
                .cornerRadius(16)
                Spacer()
            }
            
            Spacer()
        }
        .frame(maxHeight: .infinity)
        .padding()
        .background(Color("Background"))
    }
}

//#Preview {
//    SingleEvent()
//}
