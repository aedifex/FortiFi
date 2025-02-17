//
//  jwt.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/17/25.
//

import Foundation

final class JWT {
    
    static func isExpired(_ jwt: String) throws -> Bool {
        let parsed = try parseJwt(jwt)
        
        guard let expiration = parsed["exp"] as? TimeInterval else {
            throw Errors.inputError("failed to parse jwt exp time")
        }
           
        let expirationDate = Date(timeIntervalSince1970: expiration)
        return Date() >= expirationDate
    }
    
    static private func parseJwt(_ jwt: String) throws -> [String: Any] {
        let segments = jwt.components(separatedBy: ".")
        if segments.count < 3 {
            throw Errors.unauthorized("invalid jwt: \(jwt)")
        }
        return decodeJWTPartition(segments[1]) ?? [:]
    }

    static private func base64UrlDecode(_ value: String) -> Data? {
      var base64 = value
        .replacingOccurrences(of: "-", with: "+")
        .replacingOccurrences(of: "_", with: "/")

      let length = Double(base64.lengthOfBytes(using: String.Encoding.utf8))
      let requiredLength = 4 * ceil(length / 4.0)
      let paddingLength = requiredLength - length
      if paddingLength > 0 {
        let padding = "".padding(toLength: Int(paddingLength), withPad: "=", startingAt: 0)
        base64 = base64 + padding
      }
      return Data(base64Encoded: base64, options: .ignoreUnknownCharacters)
    }

    static private func decodeJWTPartition(_ value: String) -> [String: Any]? {
      guard let bodyData = base64UrlDecode(value),
        let json = try? JSONSerialization.jsonObject(with: bodyData, options: []), let payload = json as? [String: Any] else {
          return nil
      }

      return payload
    }
    
    
    
    
}
