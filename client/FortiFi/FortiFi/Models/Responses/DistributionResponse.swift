//
//  DistributionResponse.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/19/25.
//

import Foundation

struct DistributionResponse: Codable {
    var Normal: Int
    var Anomalous: Int
    var Malicious: Int
    var PrevCount: Int
}
