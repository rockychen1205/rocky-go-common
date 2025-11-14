// Package ssh provides SSH client utilities
package ssh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client represents an SSH client
type Client struct {
	client *ssh.Client
}

// Config represents SSH connection configuration
type Config struct {
	Host       string
	Port       int
	User       string
	Password   string
	KeyFile    string
	Timeout    time.Duration
}

// NewClient creates a new SSH client with the given configuration
func NewClient(config Config) (*Client, error) {
	var authMethods []ssh.AuthMethod
	
	// Add password authentication if provided
	if config.Password != "" {
		authMethods = append(authMethods, ssh.Password(config.Password))
	}
	
	// Add key-based authentication if key file provided
	if config.KeyFile != "" {
		key, err := os.ReadFile(config.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %w", err)
		}
		
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: %w", err)
		}
		
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}
	
	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication method provided")
	}
	
	// Set default timeout if not specified
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	
	clientConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: Use proper host key verification in production
		Timeout:         config.Timeout,
	}
	
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	
	return &Client{client: client}, nil
}

// Close closes the SSH connection
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// RunCommand executes a command on the remote server
func (c *Client) RunCommand(cmd string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()
	
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), fmt.Errorf("command execution failed: %w", err)
	}
	
	return string(output), nil
}

// DownloadFile downloads a file from the remote server using SCP
func (c *Client) DownloadFile(remotePath, localPath string) error {
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()
	
	// Create local file
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}
	
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()
	
	// Get remote file content
	remoteFile, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	
	if err := session.Start(fmt.Sprintf("cat %s", remotePath)); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}
	
	// Copy content
	if _, err := io.Copy(localFile, remoteFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	
	return session.Wait()
}

// UploadFile uploads a file to the remote server using SCP
func (c *Client) UploadFile(localPath, remotePath string) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer localFile.Close()
	
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()
	
	// Get file info
	fileInfo, err := localFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	
	// Create stdin pipe
	stdinPipe, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	
	// Start SCP command
	if err := session.Start(fmt.Sprintf("cat > %s", remotePath)); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}
	
	// Copy file content
	if _, err := io.CopyN(stdinPipe, localFile, fileInfo.Size()); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	
	stdinPipe.Close()
	
	return session.Wait()
}

// FileExists checks if a file exists on the remote server
func (c *Client) FileExists(remotePath string) (bool, error) {
	_, err := c.RunCommand(fmt.Sprintf("test -f %s", remotePath))
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			if exitErr.ExitStatus() == 1 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

// ListFiles lists files in a directory on the remote server
func (c *Client) ListFiles(remotePath string) ([]string, error) {
	output, err := c.RunCommand(fmt.Sprintf("ls -1 %s", remotePath))
	if err != nil {
		return nil, err
	}
	
	var files []string
	for _, line := range splitLines(output) {
		if line != "" {
			files = append(files, line)
		}
	}
	
	return files, nil
}

// Helper function to split output into lines
func splitLines(output string) []string {
	var lines []string
	var line string
	for _, ch := range output {
		if ch == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(ch)
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}
