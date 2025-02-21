//
//  Font+.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/21/25.
//

import Foundation
import SwiftUICore
import UIKit

extension View {
    func Title() -> some View {
        return self.font(.title)
            .fontWeight(.medium)
    }
    func Header() -> some View {
        self.font(.title3)
            .fontWeight(.regular)
    }
    func Label() -> some View {
        self.font(.subheadline)
    }
    func Tag() -> some View {
        self.font(.caption)
    }
}

