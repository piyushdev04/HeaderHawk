# HeaderHawk 🦅  

## What is HeaderHawk?  
HeaderHawk is a simple security header scanner built with Go. It checks if a website sends the right HTTP response headers that protect against common attacks. Think of it as a quick health check for the invisible security rules your browser should follow.  

---

## Purpose  
The goal of HeaderHawk is to make web security easier to understand and apply. Many attacks don’t happen because of broken code, but because servers are not configured safely. With HeaderHawk, you can see which important headers are missing and get advice on how to fix them.  

---

## What does it do?  
HeaderHawk:  
- Sends a safe request to your website.  
- Reads the response headers.  
- Compares them against modern security best practices.  
- Shows which headers are missing or unsafe.  
- Explains the problem in plain English and suggests a fix.  

---

## The Headers We Check  
- **Strict-Transport-Security (HSTS):** Forces HTTPS, protects from downgrade attacks.  
- **Content-Security-Policy (CSP):** Prevents XSS by locking down allowed resources.  
- **X-Frame-Options:** Stops clickjacking attacks.  
- **X-Content-Type-Options:** Prevents browsers from misinterpreting files.  
- **Referrer-Policy:** Controls what sensitive info leaks when users click links.  

---

> “I’m more concerned with being right than being fast.”  
— Gilfoyle, *Silicon Valley* 
![alt text](wp11727429-gilfoyle-wallpapers.jpg)