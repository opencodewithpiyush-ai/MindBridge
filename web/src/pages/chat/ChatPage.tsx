import { useState, useEffect, useRef } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import API_BASE from '../../api';
import './ChatPage.css';

interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  files?: { key: string; url: string }[];
}

interface Model {
  id: string;
  name: string;
  display: string;
}

interface UploadedFile {
  key: string;
  url: string;
  file: File;
}

export default function ChatPage() {
  const navigate = useNavigate();
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const [models, setModels] = useState<Model[]>([]);
  const [selectedModel, setSelectedModel] = useState<string>('');
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [attachments, setAttachments] = useState<UploadedFile[]>([]);
  
  // Auth check
  useEffect(() => {
    const token = localStorage.getItem('mb_token');
    if (!token) {
      navigate('/login');
    }
  }, [navigate]);

  // Fetch models
  useEffect(() => {
    const fetchModels = async () => {
      try {
        const res = await fetch(`${API_BASE}/models`);
        const data = await res.json();
        if (data.success && data.models) {
          setModels(data.models);
          if (data.models.length > 0) {
            setSelectedModel(data.models[0].name);
          }
        }
      } catch (err) {
        console.error('Failed to fetch models:', err);
      }
    };
    fetchModels();
  }, []);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setIsUploading(true);
    const formData = new FormData();
    formData.append('file', file);
    formData.append('name', file.name);
    formData.append('type', file.type);

    try {
      const token = localStorage.getItem('mb_token');
      const res = await fetch(`${API_BASE}/upload`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        },
        body: formData
      });

      const data = await res.json();
      if (data.success) {
        setAttachments(prev => [...prev, { key: data.key, url: data.url, file }]);
      } else {
        alert(data.error || 'Upload failed');
      }
    } catch (err) {
      console.error('Upload error:', err);
      alert('Upload failed. Please try again.');
    } finally {
      setIsUploading(false);
      if (e.target) e.target.value = ''; // reset input
    }
  };

  const removeAttachment = (index: number) => {
    setAttachments(prev => prev.filter((_, i) => i !== index));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim() && attachments.length === 0) return;
    if (isLoading) return;

    const token = localStorage.getItem('mb_token');
    if (!token) {
      navigate('/login');
      return;
    }

    const currentInput = input;
    const currentAttachments = [...attachments];
    
    // Add user message to UI
    const userMsgId = Date.now().toString();
    const newUserMsg: Message = {
      id: userMsgId,
      role: 'user',
      content: currentInput,
      files: currentAttachments.length > 0 ? currentAttachments.map(a => ({ key: a.key, url: a.url })) : undefined
    };
    
    setMessages(prev => [...prev, newUserMsg]);
    setInput('');
    setAttachments([]);
    setIsLoading(true);

    // Create a placeholder for AI message
    const aiMsgId = (Date.now() + 1).toString();
    setMessages(prev => [...prev, { id: aiMsgId, role: 'assistant', content: '' }]);

    try {
      const payload = {
        query: currentInput.trim() === '' ? ' ' : currentInput,
        model: selectedModel,
        files: currentAttachments.map(a => ({ key: a.key, url: a.url, type: a.file.type }))
      };

      const response = await fetch(`${API_BASE}/chat/stream-raw`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(payload)
      });

      if (!response.ok) {
        let errText = `Error: ${response.status}`;
        try {
          const errData = await response.json();
          if (errData.error) errText = errData.error;
        } catch (e) {}
        throw new Error(errText);
      }

      const reader = response.body?.getReader();
      const decoder = new TextDecoder('utf-8');

      if (reader) {
        let aiFullContent = '';
        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          const chunk = decoder.decode(value, { stream: true });
          const lines = chunk.split('\n');

          for (const line of lines) {
            if (line.startsWith('data: ')) {
              const dataStr = line.substring(6).trim();
              if (!dataStr) continue;
              
              try {
                const parsed = JSON.parse(dataStr);
                
                if (parsed.status === 'connected') {
                  continue;
                }
                
                if (parsed.error) {
                  throw new Error(parsed.error);
                }
                
                if (parsed.content) {
                  aiFullContent += parsed.content;
                  setMessages(prev => prev.map(msg => 
                    msg.id === aiMsgId ? { ...msg, content: aiFullContent } : msg
                  ));
                }
              } catch (e) {
                // Ignore parse errors, but if it's an explicit error thrown by us, rethrow it
                if (e instanceof Error && e.message !== 'Unexpected token' && !e.message.includes('JSON')) {
                  throw e;
                }
              }
            }
          }
        }
      }
    } catch (err: any) {
      console.error('Chat error:', err);
      const errorMessage = err.message || 'Sorry, I encountered an error. Please try again.';
      setMessages(prev => prev.map(msg => 
        msg.id === aiMsgId ? { ...msg, content: `Error: ${errorMessage}` } : msg
      ));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="chat-container">
      {/* Sidebar / Left Navigation */}
      <aside className="chat-sidebar">
        <div className="sidebar-header">
          <Link to="/" className="sidebar-logo">
            <span className="logo-emoji">🧠</span>
            <span className="logo-text">MindBridge</span>
          </Link>
        </div>

        <div className="sidebar-content">
          <button className="new-chat-btn" onClick={() => setMessages([])}>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <line x1="12" y1="5" x2="12" y2="19"></line>
              <line x1="5" y1="12" x2="19" y2="12"></line>
            </svg>
            New Chat
          </button>
        </div>

        <div className="sidebar-footer">
          <div className="model-selector-wrapper">
            <label htmlFor="model-select" className="model-label">Model</label>
            <select 
              id="model-select"
              className="model-select" 
              value={selectedModel} 
              onChange={(e) => setSelectedModel(e.target.value)}
              disabled={isLoading || models.length === 0}
            >
              {models.length === 0 && <option value="">Loading models...</option>}
              {models.map(m => (
                <option key={m.id} value={m.name}>{m.display || m.name}</option>
              ))}
            </select>
          </div>
          
          <button 
            className="logout-btn" 
            onClick={() => {
              localStorage.removeItem('mb_token');
              navigate('/login');
            }}
          >
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
              <polyline points="16 17 21 12 16 7"></polyline>
              <line x1="21" y1="12" x2="9" y2="12"></line>
            </svg>
            Logout
          </button>
        </div>
      </aside>

      {/* Main Chat Area */}
      <main className="chat-main">
        {messages.length === 0 ? (
          <div className="chat-empty-state">
            <div className="empty-icon">✨</div>
            <h2>How can I help you today?</h2>
            <p>Select a model from the sidebar and start chatting.</p>
          </div>
        ) : (
          <div className="chat-messages">
            {messages.map((msg) => (
              <div key={msg.id} className={`message-wrapper ${msg.role}`}>
                <div className="message-bubble">
                  {msg.files && msg.files.length > 0 && (
                    <div className="message-files">
                      {msg.files.map((file, i) => (
                        <img key={i} src={file.url} alt="Uploaded attachment" className="message-img" />
                      ))}
                    </div>
                  )}
                  <div className="message-content">{msg.content || (isLoading && msg.role === 'assistant' ? <span className="typing-indicator">...</span> : '')}</div>
                </div>
              </div>
            ))}
            <div ref={messagesEndRef} />
          </div>
        )}

        <div className="chat-input-container">
          <div className="chat-input-wrapper">
            
            {attachments.length > 0 && (
              <div className="attachments-preview">
                {attachments.map((att, i) => (
                  <div key={i} className="attachment-item">
                    <img src={URL.createObjectURL(att.file)} alt="preview" />
                    <button className="remove-att-btn" onClick={() => removeAttachment(i)}>
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                      </svg>
                    </button>
                  </div>
                ))}
              </div>
            )}

            <form onSubmit={handleSubmit} className="input-form">
              <label className="attach-btn" title="Upload Image">
                <input 
                  type="file" 
                  accept="image/*" 
                  onChange={handleFileChange} 
                  disabled={isUploading || isLoading}
                  style={{ display: 'none' }} 
                />
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"></path>
                </svg>
              </label>
              
              <input
                type="text"
                className="chat-input"
                placeholder={isUploading ? "Uploading..." : "Message MindBridge..."}
                value={input}
                onChange={(e) => setInput(e.target.value)}
                disabled={isLoading || isUploading}
              />
              
              <button 
                type="submit" 
                className={`send-btn ${input.trim() || attachments.length > 0 ? 'active' : ''}`}
                disabled={(!input.trim() && attachments.length === 0) || isLoading || isUploading}
              >
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <line x1="22" y1="2" x2="11" y2="13"></line>
                  <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
                </svg>
              </button>
            </form>
          </div>
          <div className="input-footer">AI models can make mistakes. Always verify important information.</div>
        </div>
      </main>
    </div>
  );
}
