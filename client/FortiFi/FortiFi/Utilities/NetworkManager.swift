//
//  NetworkManager.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/16/25.
//

import Foundation
import SwiftUI

@MainActor final class NetworkManager: Observable {
    
    static let shared = NetworkManager()
    var fcm = ""
    
    static private let baseUrl = "http://192.168.4.136:3000"
    private let loginUrl = baseUrl + "/Login"
    private let eventsUrl = baseUrl + "/GetUserEvents"
    private let refreshUrl = baseUrl + "/RefreshUser"
    private let setFcmUrl = baseUrl + "/UpdateFcm"
    private let distributionUrl = baseUrl + "/GetWeeklyDistribution"
    private let getDevicesUrl = baseUrl + "/GetDevices"
    private let getThreatAssistanceUrL = baseUrl + "/GetThreatAssistance"
    private let getMoreAssistanceUrl = baseUrl + "/GetMoreAssistance"
    private let getRecommendationsUrl = baseUrl + "/GetRecommendations"
    private let getGeneralAssistanceUrl = baseUrl + "/GetGeneralAssistance"
    
    @AppStorage("refreshToken") private var refreshToken: String = ""
    @AppStorage("jwt") var jwt: String = ""
    
    func login(_ user: User) async throws -> String {
        
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
            
            if response.value(forHTTPHeaderField: "Refresh") == nil {
                throw Errors.networkError("could not get auth tokens")
            }
            if response.value(forHTTPHeaderField: "Jwt") == nil {
                throw Errors.networkError("could not get auth tokens")
            }
            
            refreshToken = response.value(forHTTPHeaderField: "Refresh")!
            jwt = response.value(forHTTPHeaderField: "Jwt")!
            try await setNotificationsToken()
            return try JWT.getSubject(for: jwt)
            
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
        
        let requestBody = SetFcmRequest(fcm_token: fcm)
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
            throw Errors.networkError("failed to register notifications")
        case 404:
            throw Errors.notFound("no use associated with token")
        default:
            throw Errors.networkError("unable to update notifications token")
        }
        
    }
    
    private func refreshAuthTokens() async throws {
        
        guard var url = URL(string: refreshUrl) else {
            throw Errors.invalidUrl("url could not be constructed")
        }
        
        url.append(queryItems: [URLQueryItem(name: "id", value: LoginViewModel.shared.userId)])
        var request = URLRequest(url:url)
        request.setValue("\(refreshToken)", forHTTPHeaderField: "Refresh")
        request.httpMethod = "GET"
        let (_, response) = try await URLSession.shared.data(for: request)
        
        guard let response = response as? HTTPURLResponse else {
            throw Errors.internalError("failed to parse response")
        }
        
        jwt = ""
        refreshToken = ""
        switch response.statusCode {
        case 200:
            if response.value(forHTTPHeaderField: "Refresh") == nil {
                throw Errors.networkError("could not get auth tokens")
            }
            if response.value(forHTTPHeaderField: "Jwt") == nil {
                throw Errors.networkError("could not get auth tokens")
            }
            
            refreshToken = response.value(forHTTPHeaderField: "Refresh")!
            jwt = response.value(forHTTPHeaderField: "Jwt")!
        case 404:
            BaseViewModel.shared.authenticated = false
            throw Errors.notFound("user does not exist")
        case 401:
            BaseViewModel.shared.authenticated = false
            throw Errors.expiredToken("tokens have expired")
        default:
            BaseViewModel.shared.authenticated = false
            throw Errors.networkError("network error")
        }
        
    }
    
    func getEvents() async throws -> [Event] {
        
        if try JWT.isExpired(jwt) {
            try await refreshAuthTokens()
        }
        
        guard let url = URL(string: eventsUrl) else {
            throw Errors.invalidUrl("url could not be constructed")
        }
        
        var request = URLRequest(url:url)
        request.setValue("Bearer \(jwt)", forHTTPHeaderField: "Authorization")
        request.httpMethod = "GET"
        
        let (data, response) = try await URLSession.shared.data(for: request)
        
        guard let response = response as? HTTPURLResponse else {
            throw Errors.internalError("failed to parse response")
        }
        
        switch response.statusCode {
        case 200:
            let events = try JSONDecoder().decode(EventsResponse.self, from: data)
            if events.events == nil { return [] }
            return events.events!
        case 404:
            throw Errors.notFound("user does not exist")
        case 401:
            throw Errors.expiredToken("tokens have expired")
        default:
            throw Errors.networkError("network error")
        }
        
    }
    
    func getEventsDistribution() async throws -> DistributionResponse{
        
        if try JWT.isExpired(jwt) {
            try await refreshAuthTokens()
        }
        
        guard let url = URL(string: distributionUrl) else {
            throw Errors.invalidUrl("url could not be constructed")
        }
        
        var request = URLRequest(url:url)
        request.setValue("Bearer \(jwt)", forHTTPHeaderField: "Authorization")
        request.httpMethod = "GET"
        
        let (data, response) = try await URLSession.shared.data(for: request)
        
        guard let response = response as? HTTPURLResponse else {
            throw Errors.internalError("failed to parse response")
        }
        
        switch response.statusCode {
        case 200:
            let distribution = try JSONDecoder().decode(DistributionResponse.self, from: data)
            return distribution
        case 404:
            throw Errors.notFound("user does not exist")
        case 401:
            throw Errors.expiredToken("tokens have expired")
        default:
            throw Errors.networkError("network error")
        }
        
    }
    
    func getDevices() async throws -> [DevicesResponse] {
        
        if try JWT.isExpired(jwt) {
            try await refreshAuthTokens()
        }
        
        guard let url = URL(string: getDevicesUrl) else {
            throw Errors.invalidUrl("url could not be constructed")
        }
        
        var request = URLRequest(url:url)
        request.setValue("Bearer \(jwt)", forHTTPHeaderField: "Authorization")
        request.httpMethod = "GET"
        
        let (data, response) = try await URLSession.shared.data(for: request)
        
        guard let response = response as? HTTPURLResponse else {
            throw Errors.internalError("failed to parse response")
        }
        
        switch response.statusCode {
        case 200:
            let devices = try JSONDecoder().decode([DevicesResponse].self, from: data)
            return devices
        case 404:
            throw Errors.notFound("user does not exist")
        case 401:
            throw Errors.expiredToken("tokens have expired")
        default:
            throw Errors.networkError("network error")
        }
    }
    
    func getThreatAssistance(threatId: Int) async throws -> ChatMessage {
        
        if try JWT.isExpired(jwt) {
            try await refreshAuthTokens()
        }
        
        guard var url = URL(string: getThreatAssistanceUrL) else {
            throw Errors.invalidUrl("url could not be constructed")
        }
        
        url.append(queryItems: [URLQueryItem(name: "threatId", value: String(threatId))])
        var request = URLRequest(url:url)
        request.setValue("Bearer \(jwt)", forHTTPHeaderField: "Authorization")
        request.httpMethod = "GET"
        let (data, response) = try await URLSession.shared.data(for: request)
    
        guard let response = response as? HTTPURLResponse else {
            throw Errors.internalError("failed to parse response")
        }
        
        switch response.statusCode {
        case 200:
            let messageData = try JSONDecoder().decode(ChatResponse.self, from: data)
            return ChatMessage(id: messageData.id, text: messageData.response, sender: 1)
        case 401:
            throw Errors.unauthorized("invalid token")
        case 404:
            throw Errors.notFound("user or threat not found")
        case 500:
            throw Errors.internalError("internal server error")
        default:
            throw Errors.networkError("network error")
        }
        
    }
    
    func getMoreAssistance(threatId: Int, query: String) async throws -> ChatMessage {
        
        if try JWT.isExpired(jwt) {
            try await refreshAuthTokens()
        }
        
        guard var url = URL(string: getMoreAssistanceUrl) else {
            throw Errors.invalidUrl("url could not be constructed")
        }
        
        url.append(queryItems: [URLQueryItem(name: "threatId", value: String(threatId))])
        var request = URLRequest(url:url)
        
        let requestBody = ChatbotQueryRequest(query: query)
        guard let encodedData = try? JSONEncoder().encode(requestBody) else {
            throw Errors.internalError("failed to encode request body")
        }
        
        request.setValue("Bearer \(jwt)", forHTTPHeaderField: "Authorization")
        request.httpMethod = "POST"
        let (data, response) = try await URLSession.shared.upload(for: request, from:encodedData)

        guard let response = response as? HTTPURLResponse else {
            throw Errors.internalError("failed to parse response")
        }
        
        switch response.statusCode {
        case 200:
            let messageData = try JSONDecoder().decode(ChatResponse.self, from: data)
            return ChatMessage(id: messageData.id, text: messageData.response, sender: 1)
        case 401:
            throw Errors.unauthorized("invalid token")
        case 404:
            throw Errors.notFound("user or threat not found")
        case 500:
            throw Errors.internalError("internal server error")
        default:
            throw Errors.networkError("network error")
        }
    }
    
    func getRecommendations(threatId: Int) async throws -> ChatMessage{
        
        if try JWT.isExpired(jwt) {
            try await refreshAuthTokens()
        }
        
        guard var url = URL(string: getRecommendationsUrl) else {
            throw Errors.invalidUrl("url could not be constructed")
        }
        
        url.append(queryItems: [URLQueryItem(name: "threatId", value: String(threatId))])
        var request = URLRequest(url:url)
        
        request.setValue("Bearer \(jwt)", forHTTPHeaderField: "Authorization")
        request.httpMethod = "GET"
        let (data, response) = try await URLSession.shared.data(for: request)
    
        guard let response = response as? HTTPURLResponse else {
            throw Errors.internalError("failed to parse response")
        }
        
        switch response.statusCode {
        case 200:
            let messageData = try JSONDecoder().decode(ChatResponse.self, from: data)
            return ChatMessage(id: messageData.id, text: messageData.response, sender: 1)
        case 401:
            throw Errors.unauthorized("invalid token")
        case 404:
            throw Errors.notFound("user or threat not found")
        case 500:
            throw Errors.internalError("internal server error")
        default:
            throw Errors.networkError("network error")
        }
        
    }
    
    func getGeneralAssistance(query: String) async throws -> ChatMessage {
        
        if try JWT.isExpired(jwt) {
            try await refreshAuthTokens()
        }
        
        guard let url = URL(string: getGeneralAssistanceUrl) else {
            throw Errors.invalidUrl("url could not be constructed")
        }
        
        var request = URLRequest(url:url)
        
        request.setValue("Bearer \(jwt)", forHTTPHeaderField: "Authorization")
        request.httpMethod = "POST"
        
        let requestBody = ChatbotQueryRequest(query: query)
        guard let encodedData = try? JSONEncoder().encode(requestBody) else {
            throw Errors.internalError("failed to encode request body")
        }
        
        let (data, response) = try await URLSession.shared.upload(for: request, from: encodedData)
    
        guard let response = response as? HTTPURLResponse else {
            throw Errors.internalError("failed to parse response")
        }
        
        switch response.statusCode {
        case 200:
            let messageData = try JSONDecoder().decode(ChatResponse.self, from: data)
            return ChatMessage(id: messageData.id, text: messageData.response, sender: 1)
        case 401:
            throw Errors.unauthorized("invalid token")
        case 404:
            throw Errors.notFound("user or threat not found")
        case 500:
            throw Errors.internalError("internal server error")
        default:
            throw Errors.networkError("network error")
        }
        
    }
}

