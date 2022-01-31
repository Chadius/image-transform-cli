# Why
Applications need to recognize errors when the package crashes.

# What is it
- CLI should print a failure when it uses the local library
- CLI should print out 500 failure from Server

# What happens next
- Acts as a blueprint for later clients

# Caveats
Server should NOT return the error message. Instead, we should log the error.
We shouldn't make our errors visible to clients using the web interface - they likely won't care, and we shouldn't expose our implementation later.
