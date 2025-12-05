#!/usr/bin/env python3
"""
Auth Service Test Script
Tests login, re-authentication, RBAC, and service accounts
"""

import json
import requests
import time
from datetime import datetime

BASE_URL = "http://localhost:8080"

class Colors:
    GREEN = '\033[0;32m'
    BLUE = '\033[0;34m'
    RED = '\033[0;31m'
    YELLOW = '\033[1;33m'
    NC = '\033[0m'

def print_header(text):
    print(f"\n{Colors.BLUE}{'='*60}{Colors.NC}")
    print(f"{Colors.BLUE}{text}{Colors.NC}")
    print(f"{Colors.BLUE}{'='*60}{Colors.NC}")

def print_success(text):
    print(f"{Colors.GREEN}✅ {text}{Colors.NC}")

def print_error(text):
    print(f"{Colors.RED}❌ {text}{Colors.NC}")

def print_info(text):
    print(f"{Colors.YELLOW}ℹ️  {text}{Colors.NC}")

def test_login(username, password):
    """Test login endpoint"""
    print_header(f"Testing Login: {username}")
    
    try:
        response = requests.post(
            f"{BASE_URL}/api/v1/login",
            json={"username": username, "password": password}
        )
        
        if response.status_code == 200:
            data = response.json()
            print_success(f"Login successful")
            print(f"   Token: {data['token'][:50]}...")
            return data['token']
        else:
            print_error(f"Login failed: {response.status_code}")
            print(f"   Response: {response.text}")
            return None
    except Exception as e:
        print_error(f"Login error: {e}")
        return None

def test_re_auth(token, password):
    """Test re-authentication endpoint (FDA 21 CFR Part 11)"""
    print_header("Testing Re-Authentication")
    
    try:
        response = requests.post(
            f"{BASE_URL}/api/v1/re-auth",
            json={"password": password},
            headers={"Authorization": f"Bearer {token}"}
        )
        
        if response.status_code == 200:
            data = response.json()
            print_success("Re-authentication successful")
            print(f"   Signing Token: {data['signing_token'][:50]}...")
            return data['signing_token']
        else:
            print_error(f"Re-auth failed: {response.status_code}")
            print(f"   Response: {response.text}")
            return None
    except Exception as e:
        print_error(f"Re-auth error: {e}")
        return None

def test_create_service_account(admin_token, name):
    """Test service account creation (ADMIN only)"""
    print_header(f"Testing Service Account Creation: {name}")
    
    try:
        response = requests.post(
            f"{BASE_URL}/api/v1/service-accounts",
            json={"name": name},
            headers={"Authorization": f"Bearer {admin_token}"}
        )
        
        if response.status_code == 200:
            data = response.json()
            print_success("Service account created")
            print(f"   Token: {data['token'][:50]}...")
            print(f"   Expires: Never (long-lived)")
            return data['token']
        else:
            print_error(f"Service account creation failed: {response.status_code}")
            print(f"   Response: {response.text}")
            return None
    except Exception as e:
        print_error(f"Service account error: {e}")
        return None

def test_rbac_protection(token, role):
    """Test RBAC middleware - try to access admin endpoint"""
    print_header(f"Testing RBAC Protection ({role} role)")
    
    try:
        response = requests.post(
            f"{BASE_URL}/api/v1/service-accounts",
            json={"name": "test-service"},
            headers={"Authorization": f"Bearer {token}"}
        )
        
        if role == "ADMIN":
            if response.status_code == 200:
                print_success("ADMIN can access protected endpoint")
            else:
                print_error(f"ADMIN access denied: {response.status_code}")
        else:
            if response.status_code == 403:
                print_success(f"{role} correctly denied access (403)")
            else:
                print_error(f"Expected 403, got {response.status_code}")
                
    except Exception as e:
        print_error(f"RBAC test error: {e}")

def test_invalid_token():
    """Test with invalid token"""
    print_header("Testing Invalid Token")
    
    try:
        response = requests.post(
            f"{BASE_URL}/api/v1/service-accounts",
            json={"name": "test"},
            headers={"Authorization": "Bearer invalid-token-12345"}
        )
        
        if response.status_code == 401:
            print_success("Invalid token correctly rejected (401)")
        else:
            print_error(f"Expected 401, got {response.status_code}")
    except Exception as e:
        print_error(f"Invalid token test error: {e}")

def test_missing_token():
    """Test without token"""
    print_header("Testing Missing Token")
    
    try:
        response = requests.post(
            f"{BASE_URL}/api/v1/service-accounts",
            json={"name": "test"}
        )
        
        if response.status_code == 401:
            print_success("Missing token correctly rejected (401)")
        else:
            print_error(f"Expected 401, got {response.status_code}")
    except Exception as e:
        print_error(f"Missing token test error: {e}")

def main():
    print(f"\n{Colors.BLUE}🔐 Auth Service Test Suite{Colors.NC}")
    print(f"{Colors.BLUE}{'='*60}{Colors.NC}")
    print(f"Base URL: {BASE_URL}")
    print(f"Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    
    # Test 1: Login as admin
    admin_token = test_login("admin", "admin123")
    if not admin_token:
        print_error("Cannot continue without admin token")
        return
    
    time.sleep(1)
    
    # Test 2: Re-authentication
    new_admin_token = test_re_auth(admin_token, "admin123")
    if new_admin_token:
        admin_token = new_admin_token
    
    time.sleep(1)
    
    # Test 3: Create service account (ADMIN only)
    service_token = test_create_service_account(admin_token, "test-service-1")
    
    time.sleep(1)
    
    # Test 4: RBAC - Admin can access
    test_rbac_protection(admin_token, "ADMIN")
    
    time.sleep(1)
    
    # Test 5: RBAC - Service account cannot access admin endpoint
    if service_token:
        test_rbac_protection(service_token, "SERVICE")
    
    time.sleep(1)
    
    # Test 6: Invalid token
    test_invalid_token()
    
    time.sleep(1)
    
    # Test 7: Missing token
    test_missing_token()
    
    # Summary
    print_header("Test Summary")
    print_success("All basic tests completed!")
    print_info("Check the logs above for detailed results")
    
    print(f"\n{Colors.YELLOW}💡 Useful Commands:{Colors.NC}")
    print("  - View logs: docker-compose logs -f auth")
    print("  - Check database: docker exec ops-postgres-1 psql -U postgres -d historian -c 'SELECT * FROM users;'")
    print("  - Stop services: cd ops && docker-compose down")

if __name__ == "__main__":
    main()
