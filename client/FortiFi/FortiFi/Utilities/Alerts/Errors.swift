//
//  Errors.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/16/25.
//

import Foundation

enum Errors: Error {
    case invalidUrl(String), networkError(String), inputError(String), notFound(String),
    unauthorized(String), internalError(String), expiredToken(String)
}
