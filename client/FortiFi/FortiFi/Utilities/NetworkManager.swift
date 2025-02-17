//
//  NetworkManager.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/16/25.
//

import Foundation
import SwiftUI

final class NetworkManager {
    
    static let shared = NetworkManager()
    static var fcm = ""

    static private let baseUrl = "<fill in with ip address:port of device the server is running on"
    private let loginUrl = baseUrl + "/Login"
    private let eventsUrl = baseUrl + "/GetUserEvents"
    private let refreshUrl = baseUrl + "/RefreshUser"
    private let setFcmUrl = baseUrl + "/UpdateFcm"
    private var jwt = ""
        
    @AppStorage("refreshToken") private var refreshToken: String = ""
    
    func login(_ user: User) async throws {

        guard let url = URL(string: loginUrl) else {
            throw Errors.invalidUrl("url could not be constructed")
        }
        
        let requestBody = LoginRequest(user: user)
        guard let encodedData = try? JSONEncoder().encode(requestBody) else {
            throw Errors.inputError("invalid inputs")
        }
        
        var request = URLRequest(url:url)
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpMethod = "POST"
        
        let (_, response) = try await URLSession.shared.upload(for: request, from:encodedData)
        
        guard let response = response as? HTTPURLResponse else {
            throw Errors.internalError("failed to parse response")
        }

        switch response.statusCode {
        case 200:
            let headers = response.allHeaderFields
            guard let refreshHeader = headers["Refresh"] as? String else {
                throw Errors.networkError("could not get auth tokens")
            }
            guard let jwtHeader = headers["Jwt"] as? String else {
                throw Errors.networkError("could not get auth tokens")
            }
            
            refreshToken = refreshHeader
            jwt = jwtHeader
            try await setNotificationsToken()
        case 404:
            throw Errors.notFound("user does not exist")
        case 401:
            throw Errors.unauthorized("invalid password")
        case 400:
            throw Errors.inputError("invalid inputs")
        default:
            throw Errors.networkError("network error")
        }
            
    }
    
    private func setNotificationsToken() async throws {
        
        // check token
        if try JWT.isExpired(jwt) {
            try await refreshAuthTokens()
        }
        
        guard let url = URL(string: setFcmUrl) else {
            throw Errors.invalidUrl("url could not be constructed")
        }
        
        let requestBody = SetFcmRequest(fcm_token: NetworkManager.fcm)
        guard let encodedData = try? JSONEncoder().encode(requestBody) else {
            throw Errors.internalError("failed to encode request body")
        }
        
        var request = URLRequest(url:url)
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(jwt)", forHTTPHeaderField: "Authorization")
        request.httpMethod = "POST"
        
        let (_, response) = try await URLSession.shared.upload(for: request, from:encodedData)
        guard let response = response as? HTTPURLResponse else {
            throw Errors.internalError("failed to parse response")
        }
        switch response.statusCode {
        case 202:
            break
        case 401:
            throw Errors.unauthorized("invalid auth token")
        case 400:
            throw Errors.inputError("invalid fcm token")
        case 404:
            throw Errors.notFound("no use associated with token")
        default:
            throw Errors.networkError("unable to update notifications token")
        }
        
    }
    
    private func refreshAuthTokens() async throws {
        
        guard let url = URL(string: refreshUrl) else {
            throw Errors.invalidUrl("url could not be constructed")
        }
        
        var request = URLRequest(url:url)
        request.setValue("\(refreshToken)", forHTTPHeaderField: "Refresh")
        request.httpMethod = "GET"
        
        let (_, response) = try await URLSession.shared.upload(for: request, from: Data())
        
        guard let response = response as? HTTPURLResponse else {
            throw Errors.internalError("failed to parse response")
        }
        
        switch response.statusCode {
        case 200:
            let headers = response.allHeaderFields
            guard let refreshHeader = headers["Refresh"] as? String else {
                throw Errors.networkError("could not get auth tokens")
            }
            guard let jwtHeader = headers["Jwt"] as? String else {
                throw Errors.networkError("could not get auth tokens")
            }
            refreshToken = refreshHeader
            jwt = jwtHeader
        case 404:
            throw Errors.notFound("user does not exist")
        case 401:
            throw Errors.expiredToken("tokens have expired")
        default:
            throw Errors.networkError("network error")
        }
        
    }
    
    
}

