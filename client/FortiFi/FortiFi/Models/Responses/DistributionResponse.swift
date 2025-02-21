//
//  DistributionResponse.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/19/25.
//

import Foundation

struct DistributionResponse: Codable {
    var Benign: Int
    var PortScan: Int
    var DDoS: Int
    var PrevWeekTotal: Int
}
