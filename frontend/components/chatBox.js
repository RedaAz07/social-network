"use client";
import { useState, useRef, useEffect } from "react";
import AddReactionIcon from "@mui/icons-material/AddReaction";
import SendIcon from "@mui/icons-material/Send";
import Link from "next/link";
import "../styles/chat.css";
import { useWS } from "../context/wsContext.js";
import { useParams } from "next/navigation";
import { useProfile } from "../context/profile.js";

export default function ChatBox({ user }) {
  const { sendMessage, addListener, removeListener } = useWS();
  const { Profile } = useProfile();
  const [messages, setMessages] = useState([]);
  const [preview, setPreview] = useState(null);
  const [input, setInput] = useState("");
  const [showEmojis, setShowEmojis] = useState(false);
  const [toast, setToast] = useState(null)
  const [onlineUsers, setonlineUsers] = useState([])
  const [New, setNew] = useState(false)
  const inputRef = useRef(null);
  const chatEndRef = useRef(null);
  const chatBoxRef = useRef(null);
  const refscroll = useRef(0)
  const id = useParams().id;

  const showToast = (message, type = "error", duration = 3000) => {
    setToast({ message, type });
    setTimeout(() => {
      setToast(null);
    }, duration);
  };
  const scrollToBottom = () => {
    if (chatEndRef.current) {
      chatEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  };
  const handleScrollDownClick = () => {
    scrollToBottom();
    setNew(false);
  };
  const handleScroll = () => {
    const chatBox = chatBoxRef.current;
    if (!chatBox) return;
    const atBottom = chatBox.scrollHeight - chatBox.scrollTop - chatBox.clientHeight < 50;

    if (atBottom) {
      setNew(false);
    } else {
      setNew(true);
    }
  };


  if (id == "0" || !id) {
    return (
      <div className="no-chat-selected">
        <h2>No chat selected</h2>
      </div>
    );
  }
  setTimeout(() => {
    inputRef.current?.focus();
  }, 0);

  useEffect(() => {
    refscroll.current = setTimeout(() => handleScrollDownClick(), 100);
    return (() => clearTimeout(refscroll.current))
  }, []);

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (!file) {

      setPreview(null);
      return;
    }
    console.log(file.type,file.size);
    

    const allowedTypes = ["image/png", "image/jpeg", "image/jpg", "image/gif", "image/webp"];
    if (!allowedTypes.includes(file.type )|| file.size>= 1024*1024*30 ) {
      showToast("Only image files are allowed! || big size")
      e.target.value = "";
      setPreview(null);
      return;
    }

   const reader = new FileReader();

reader.onload = () => {
  const arr = new Uint8Array(reader.result);

  const isPng  = arr[0] === 0x89 && arr[1] === 0x50 && arr[2] === 0x4E && arr[3] === 0x47;
  const isJpeg = arr[0] === 0xFF && arr[1] === 0xD8 && arr[2] === 0xFF;
  const isGif  = arr[0] === 0x47 && arr[1] === 0x49 && arr[2] === 0x46;
  const isWebP = arr[0] === 0x52 && arr[1] === 0x49 && arr[2] === 0x46 && arr[3] === 0x46 && arr[8] === 0x57 && arr[9] === 0x45 && arr[10] === 0x42 && arr[11] === 0x50;
  const isBmp  = arr[0] === 0x42 && arr[1] === 0x4D;
  const isTiff = (arr[0] === 0x49 && arr[1] === 0x49 && arr[2] === 0x2A && arr[3] === 0x00) ||
                 (arr[0] === 0x4D && arr[1] === 0x4D && arr[2] === 0x00 && arr[3] === 0x2A);

  if (!isPng && !isJpeg && !isGif && !isWebP && !isBmp && !isTiff) {
    showToast('is not a image ')
    e.target.value = "";
    setPreview(null);
    return;
  }

  const previewReader = new FileReader();
  previewReader.onloadend = () => setPreview(previewReader.result);
  previewReader.readAsDataURL(file);
};

reader.onerror = () => {
  showToast('erroor in image')
  e.target.value = "";
  setPreview(null);
};

reader.readAsArrayBuffer(file); 
  };



  const removeImage = () => {
    setPreview(null);
    document.getElementById("idfile").value = "";
  };

  // listen for online/offline status updates
  useEffect(() => {
    const handleOlineUser = (data) => {
      setonlineUsers(data.users)
    }
    sendMessage({ type: "online_list" })
    addListener("online_list", handleOlineUser)
    return () => removeListener("online_list", handleOlineUser)
  }, [addListener, removeListener])

  useEffect(() => {
    const handleLogout = (data) => {
      let useroff = data.userID
      let arr = onlineUsers.filter((id) => {
        return id !== useroff
      })
      setonlineUsers([...arr])
    }

    addListener("logout", handleLogout)
    return () => removeListener("logout", handleLogout)
  }, [addListener, removeListener])




  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const response = await fetch(
          `http://localhost:8080/api/getmessages?receiverId=${user.id}`,
          { credentials: "include", method: "GET" }
        );
        if (!response.ok) throw new Error("Failed to fetch messages");
        const data = await response.json();
        if (data.messages) {
          const formattedMessages = data.messages
            .map((msg) => ({
              text: msg.content,
              sender: msg.senderId === user.id ? "them" : "me",
              time: new Date(msg.createdAt).toLocaleString(),
              name: msg.first_name + " " + msg.last_name,
              image: msg.photo,
              PictureSend: msg.PictureSend
            }))
            .reverse();

          setMessages(formattedMessages);
        }
      } catch (error) {
        console.error("Error fetching messages:", error);
      }
    };
    fetchMessages();
  }, [user.id]);

  useEffect(() => {
    const handleIncomingMessage = (data) => {
      console.log("New message received in chatBox:", data);
      if (data.from === user.id || data.to === user.id) {
        setMessages((prev) => [
          ...prev,
          {
            text: data.content,
            sender: data.from === user.id ? "them" : "me",
            time: data.time,
            name: data.name,
            image: data.image,
            PictureSend: data.PictureSend
          },
        ]);
        handleScroll()
      }
    };

    addListener("message", handleIncomingMessage);
    return () => removeListener("message", handleIncomingMessage);
  }, [addListener, removeListener]);



  const emojiArray = ["😀", "😃", "😄", "😁", "😆", "😅", "🤣", "😂", "🚀", "💡", "😊", "😇", "🙂", "🙃", "😉", "😍", "🥰", "😘", "😗", "😋", "😛", "😜", "🤪", "😝", "🤑", "🤗", "🤭", "🤔", "🤨", "😐", "😑", "😶", "😏", "😒", "🙄", "😬", "😔", "😪", "🤤", "😴", "😷", "🤒", "🤕", "🤢", "🤮", "🥴", "😵", "🤯", "😎", "🤓"];
  const handleSendMessage = () => {
    if (input.trim() === "" && !preview) return;
    if (input.length > 1000) {
      showToast("message is too long")
      return
    }
    
    const payload = {
      receiverId: user.id,
      messageContent: input,
      type: "message",
      PictureSend: preview
    };
    sendMessage(payload);
    setInput("");
    removeImage()
    setShowEmojis(false);
    setTimeout(() => handleScrollDownClick(), 100);
  };

  const addEmoji = (emoji) => {
    const cursorPos = inputRef.current.selectionStart;
    const newText = input.slice(0, cursorPos) + emoji + input.slice(cursorPos);
    setInput(newText);
    setTimeout(() => {
      inputRef.current.focus();
      inputRef.current.setSelectionRange(cursorPos+emoji.length, cursorPos+emoji.length);
    }, 0);
  };

  return (
    <div className="chat-container">
      <div className="chat-header">
        {toast && (
          <div className={`toast ${toast.type}`}>
            <span>{toast.message}</span>
            <button onClick={() => setToast(null)} className="toast-close">×</button>
          </div>
        )}
        <Link href={`/profile/${user.id}`}>
          <img
            src={
              user?.image ? `/uploads/${user.image}` : "/assets/default.png"
            }
            alt="user avatar"
          />
        </Link>
        <div className="onlinee">
          <Link href={`/profile/${user.id}`}>
            <span className="username">{user.first_name + " " + user.last_name}</span>
          </Link>
          <span className={onlineUsers.includes(user.id) ? "on" : "off"}>
            {onlineUsers.includes(user.id) ? "online" : "offline"}
          </span>
        </div>
      </div>

      {/* Chat box */}
      <div className="chat-box" ref={chatBoxRef} onScroll={handleScroll}>
        {messages.length === 0 ? (
          <p className="no-msg">No messages yet</p>
        ) : (
          messages.map((msg, index) => (
            <div key={index} className={`message ${msg.sender}`}>
              <div className="message-header">
                <img src={msg?.image ? `/uploads/${msg?.image}` : "/assets/default.png"} alt="profile picture" className="profilePic" />
                <span className="message-name">{msg.name}</span>
              </div>
              <div className="message-bubble">
                {msg.PictureSend && (
                  <img
                    src={`/uploads/${msg.PictureSend}`}
                    alt="message image"
                    className="message-img"
                  />
                )}
                <div className="msg-content">{msg.text}</div>
              </div>
              <div className="info-time">
                <span className="time">
                  {new Date(msg.time).toLocaleTimeString("en-US", {
                    hour: 'numeric',
                    minute: '2-digit'
                  })}
                </span>
              </div>
            </div>
          ))
        )}
        {preview && (
          <div className="preview-container">
            <img src={preview} alt="preview" className="preview-image" />
            <button className="remove-btn" onClick={removeImage}>
              ×
            </button>
          </div>
        )}
        {New && (
          <div className="scroll-down-icon" onClick={handleScrollDownClick}>
            <i className="fa-solid fa-arrow-down"></i>
          </div>
        )}

        <div ref={chatEndRef}></div>
      </div>

      {/* Emoji panel */}
      {
        showEmojis && (
          <div className="emoji-panel">
            {emojiArray.map((emoji, i) => (
              <span key={i} className="emoji" onClick={() => addEmoji(emoji)}>
                {emoji}
              </span>
            ))}
          </div>
        )
      }

      {/* Input area */}
      <div className="input-area">
        <button
          className="emoji-toggle"
          onClick={() => setShowEmojis(!showEmojis)}
        >
          <AddReactionIcon />
        </button>

        <input
          ref={inputRef}
          type="text"
          placeholder="Type a message..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSendMessage()}
        />
        <input
          hidden
          type="file"
          id="idfile"
          accept="image/*"
          onChange={handleFileChange}
        />
        <label htmlFor="idfile">
          <div className="emoji-toggle">
            <i className="fa-solid fa-image " ></i>
          </div>
        </label>
        <button onClick={handleSendMessage}>
          <SendIcon />
        </button>
      </div>
    </div >
  );
}
