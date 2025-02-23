//
//  Devices.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/19/25.
//

import SwiftUI

struct Devices: View {
    @ObservedObject var viewModel = DevicesViewModel()
    var body: some View {
        NavigationStack{
            VStack (alignment: .leading, spacing: 24){
                Text("Devices")
                    .Title()
                    .fontWeight(.medium)
                    .frame(maxWidth: .infinity, alignment: .leading)
                Text("\(viewModel.devices.count) Monitored Devices")
                    .Label()
                    .foregroundColor(.foregroundMuted)
                ScrollView{
                    VStack {
                        ForEach(viewModel.devices) {device in
                            DeviceListItem(device: device)
                            if device.self != viewModel.devices.last.self {
                                       Divider()
                            }
                        }
                    }
                    .padding()
                    .background(.fortifiBackground)
                    .cornerRadius(16)
                    .shadow(color: Color.black.opacity(0.1), radius: 3, x: 0, y: 0)
                }
                .contentMargins(2)

            }
            .padding()
            .toolbarBackground(.fortifiBackground, for: .navigationBar)
            .frame(maxHeight: .infinity)
            .background(.backgroundAlt)
            .foregroundStyle(.fortifiForeground)
        }
        .onAppear{
            Task {
                await viewModel.refresh()
            }
        }
        .refreshable {
            Task {
                await viewModel.refresh()
            }
        }
    }
}

struct DeviceListItem: View {
    var device: DevicesResponse
    var body: some View {
        NavigationLink(destination: DeviceInfo(device: device)) {
            HStack {
                Text(device.name)
                    .Label()
                Spacer()
                Image(systemName: "chevron.right")
                    .foregroundStyle(.foregroundMuted)
            }
            .padding(.horizontal)
            .padding(.vertical, 16)
            .background(.fortifiBackground)
        }
    }
}

#Preview {
    Devices()
}
