"use client"
import { useState, useEffect, useRef } from "react"
import { useRouter } from 'next/navigation';
import "../styles/comment.css"


export default function Comment({ comments, isOpen, onClose, postId, onCommentChange, lodinggg, ongetcomment, post, showToast, ID }) {
  const [commentContent, setCommentContent] = useState("")
  const [selectedFile, setSelectedFile] = useState(null)
  const [loading, setLoading] = useState(false)
  const [showEmojis, setShowEmojis] = useState(false)
  const [activeCategory, setActiveCategory] = useState("smileys")
  const modalRef = useRef(null)
  const commentsContainerRef = useRef(null)
  const [scrollPos, setScrollPos] = useState(0)
  const commentRefs = useRef({})
  const fileInputRef = useRef(null)
  const emojiRef = useRef(null)
  const textareaRef = useRef(null)
  const router = useRouter()




  // Emojis organisés par catégories
  const emojiCategories = {
    smileys: [
      "😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣", "😊", "😇",
      "🙂", "🙃", "😉", "😌", "😍", "🥰", "😘", "😗", "😙", "😚",
      "😋", "😛", "😝", "😜", "🤪", "🤨", "🧐", "🤓", "😎", "🤩",
      "🥳", "😏", "😒", "😞", "😔", "😟", "😕", "🙁", "☹️", "😣",
      "😖", "😫", "😩", "🥺", "😢", "😭", "😤", "😠", "😡", "🤬"
    ],
    people: [
      "👶", "🧒", "👦", "👧", "🧑", "👨", "👩", "🧓", "👴", "👵",
      "🙍", "🙎", "🙅", "🙆", "💁", "🙋", "🧏", "🙇", "🤦", "🤷",
      "👮", "💂", "👷", "🤴", "👸", "👳", "👲", "🧕", "🤵", "👰",
      "🤰", "🤱", "👼", "🎅", "🤶", "🦸", "🦹", "🧙", "🧚", "🧛"
    ],
    animals: [
      "🐵", "🐒", "🦍", "🐶", "🐕", "🐩", "🐺", "🦊", "🐱", "🐈",
      "🦁", "🐯", "🐅", "🐆", "🐴", "🐎", "🦄", "🦓", "🦌", "🐮",
      "🐂", "🐃", "🐄", "🐷", "🐖", "🐗", "🐽", "🐏", "🐑", "🐐",
      "🐪", "🐫", "🦙", "🦒", "🐘", "🦏", "🦛", "🐭", "🐁", "🐀"
    ],
    food: [
      "🍎", "🍐", "🍊", "🍋", "🍌", "🍉", "🍇", "🍓", "🫐", "🍈",
      "🍒", "🍑", "🥭", "🍍", "🥥", "🥝", "🍅", "🍆", "🥑", "🥦",
      "🥬", "🥒", "🌶", "🫑", "🌽", "🥕", "🫒", "🧄", "🧅", "🥔",
      "🍠", "🥐", "🥯", "🍞", "🥖", "🥨", "🧀", "🥚", "🍳", "🧈"
    ],
    activities: [
      "⚽", "🏀", "🏈", "⚾", "🥎", "🎾", "🏐", "🏉", "🥏", "🎱",
      "🪀", "🏓", "🏸", "🏒", "🏑", "🥍", "🏏", "🎿", "⛷", "🏂",
      "🪂", "🏋️", "🤼", "🤸", "⛹️", "🤾", "🏌️", "🏇", "🧘", "🏄"
    ],
    travel: [
      "🚗", "🚕", "🚙", "🚌", "🚎", "🏎", "🚓", "🚑", "🚒", "🚐",
      "🛻", "🚚", "🚛", "🚜", "🏍", "🛵", "🚲", "🛴", "🛹", "🛼",
      "🚁", "✈️", "🛩", "🛫", "🛬", "🚀", "🛸", "🚂", "🚊", "🚉"
    ],
    objects: [
      "💡", "🔦", "🏮", "🪔", "📔", "📕", "📖", "📗", "📘", "📙",
      "📚", "📓", "📒", "📃", "📜", "📄", "📰", "🗞", "📑", "🔖",
      "🏷", "💰", "🪙", "💴", "💵", "💶", "💷", "💸", "💳", "🧾"
    ],
    symbols: [
      "❤️", "🧡", "💛", "💚", "💙", "💜", "🖤", "🤍", "🤎", "💔",
      "❣️", "💕", "💞", "💓", "💗", "💖", "💘", "💝", "💟", "☮️",
      "✝️", "☪️", "🕉", "☸️", "✡️", "🔯", "🕎", "☯️", "☦️", "🛐"
    ]
  }

  const categoryIcons = {
    smileys: "😀",
    people: "👨",
    animals: "🐶",
    food: "🍎",
    activities: "⚽",
    travel: "🚗",
    objects: "💡",
    symbols: "❤️"
  }

  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.scrollTop = textareaRef.current.scrollHeight;
    }
  }, [commentContent]);

  function scrollToComment(commentId) {
    const el = commentRefs.current[commentId]
    if (el) {
      el.scrollIntoView({ behavior: "smooth", block: "start" })
    }
  }

  // Fermer le picker d'emojis quand on clique ailleurs
  useEffect(() => {
    function handleClickOutside(event) {
      if (emojiRef.current && !emojiRef.current.contains(event.target)) {
        setShowEmojis(false)
      }
    }

    document.addEventListener("mousedown", handleClickOutside)
    return () => {
      document.removeEventListener("mousedown", handleClickOutside)
    }
  }, [])



  useEffect(() => {
    if (!commentsContainerRef.current) return

    const commentsContainer = commentsContainerRef.current
    const reachedBottom =
      commentsContainer.scrollTop + commentsContainer.clientHeight >= commentsContainer.scrollHeight - 5

    async function getcomment() {
      const r = await ongetcomment(post)
      if (r) {
        scrollToComment(r)
      }
    }

    if (reachedBottom && !lodinggg) {
      getcomment()
    }
  }, [scrollPos])

  const handleFileSelect = (e) => {
   if (e.target.files[0]?.size<= 2 * 1024 * 1024){

      if (e.target.files[0]?.type){
        if (e.target.files[0]?.type.split('/')[0] =='image' )
          setSelectedFile(e.target.files[0]);
      }
    }else{
      setSelectedFile(null)
      showToast('file too large')
    }
  }

  const handleRemoveFile = () => {
    setSelectedFile(null)
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  const addEmoji = (emoji) => {
    console.log(emoji.length);

    const cursorPos = textareaRef.current.selectionStart;
    const newText = commentContent.slice(0, cursorPos) + emoji + commentContent.slice(cursorPos);
    setCommentContent(newText);
    setTimeout(() => {
      textareaRef.current.focus();
      textareaRef.current.setSelectionRange(cursorPos + 2, cursorPos + 2);
    }, 0);
  }

  const handleTextChange = (e) => {
    setCommentContent(e.target.value);
  }

  async function handlePostComment(e) {
    e.preventDefault()
    if (!commentContent.trim() && !selectedFile) {
      showToast("Please enter a comment or select a file.")
      return
    }

    try {
      setLoading(true)
      const formData = new FormData()
      formData.append("postId", postId)
      formData.append("content", commentContent)
      if (selectedFile) {
        formData.append("media", selectedFile)
      }
      formData.append("whatis", window.location.href.split('/')[3])

      formData.append("groupId", ID)

      const response = await fetch(`http://localhost:8080/api/createcomment`, {
        method: "POST",
        credentials: "include",
        body: formData,
      })


      const res = await response.json()
      if (res.error) {
        if (res.error == "Unauthorized") {
          router.push("/login");
          return
        } else {
          setCommentContent("")
          setSelectedFile(null)
          showToast(res.error)
          return
        }
      }


      setCommentContent("")
      setSelectedFile(null)

      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }

      if (onCommentChange) {
        onCommentChange(res.comment_id)
        commentsContainerRef.current.scrollTop = 0
      }
    } catch (err) {

    } finally {
      setLoading(false)
    }
  }

  if (!isOpen) return null

  return (
    <div className={`cm-overlay ${isOpen ? "cm-active" : ""}`} onClick={onClose}>
      <div className="cm-modal" onClick={(e) => e.stopPropagation()} ref={modalRef}>
        <div className="cm-header">
          <button className="cm-close" aria-label="Close modal" onClick={onClose}>
            <span className="icon-close">×</span>
          </button>
          <h3 className="cm-heading">Comments</h3>
        </div>

        <div className="cm-content">
          <div
            className="cm-list-wrap"
            ref={commentsContainerRef}
            onScroll={(e) => setScrollPos(e.target.scrollTop)}
          >
            {comments && comments.length > 0 ? (
              <div className="cm-list">
                {comments.map((comment) => (
                  <div
                    key={comment.id}
                    id={`comment-${comment.id}`}
                    className="cm-item"
                    ref={(el) => (commentRefs.current[comment.id] = el)}
                  >
                    <div className="cm-avatar cm-avatar-sm">{comment.first_name?.charAt(0) || "U"}</div>
                    <div className="cm-text">
                      <div className="cm-meta">
                        <span className="cm-name">{comment.first_name + " " + comment.last_name}</span>
                        <span className="cm-time">
                          {new Date(comment.created_at).toLocaleTimeString([], {
                            hour: "2-digit",
                            minute: "2-digit",
                          })}
                        </span>
                      </div>
                      <p className="cm-msg">{comment.content}</p>

                      {comment.media_path && (
                        comment.media_path.endsWith(".mp4") || comment.media_path.endsWith(".webm") ? (
                          <video controls className="cm-media" loop autoPlay>
                            <source src={`../${comment.media_path}`} type="video/mp4" />
                            Your browser does not support the video tag.
                          </video>
                        ) : (
                          <img src={`../${comment.media_path}`} alt="Comment Media" className="cm-media" />
                        )
                      )}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="cm-empty">
                <div className="icon-chat">💬</div>
                <p className="cm-empty-title">No comments yet</p>
                <p className="cm-empty-text">Start the conversation</p>
              </div>
            )}
          </div>

          <div className="cm-form-wrap">
            <form onSubmit={handlePostComment} className="cm-form" encType="multipart/form-data">
              <div className={`cm-input-wrap ${commentContent.trim() ? 'has-text' : ''}`}>
                <textarea
                  ref={textareaRef}
                  className="cm-input"
                  placeholder="Add a comment..."
                  value={commentContent}
                  onChange={handleTextChange}
                  onInput={(e) => {
                    e.target.scrollTop = e.target.scrollHeight;
                  }}
                  rows={1}
                  required={!selectedFile}
                  disabled={loading}
                  autoFocus={true}
                  onBlur={() => textareaRef.current?.focus()}

                />

                {/* Bouton emoji */}
                <button
                  type="button"
                  className={`cm-emoji-btn ${showEmojis ? 'active' : ''}`}
                  onClick={() => setShowEmojis(!showEmojis)}
                >
                  😊
                </button>

                {/* Input file */}
                <input
                  type="file"
                  accept="image/*,video/*"
                  onChange={handleFileSelect}
                  className="cm-file"
                  id="comment-file"
                  ref={fileInputRef}
                />
                <label
                  htmlFor="comment-file"
                  className={`cm-file-label ${selectedFile ? 'has-file' : ''}`}
                >
                  {selectedFile && <span className="cm-file-selected"></span>}
                </label>

                <button
                  type="submit"
                  className="cm-btn"
                  disabled={loading || (!commentContent.trim() && !selectedFile)}
                >
                  {loading ? "..." : "Post"}
                </button>
              </div>

              {/* Picker d'emojis avec catégories */}
              {showEmojis && (
                <div className="cm-emoji-picker" ref={emojiRef}>
                  <div className="cm-emoji-header">
                    <span>Emojis</span>
                    <button
                      type="button"
                      className="cm-emoji-close"
                      onClick={() => setShowEmojis(false)}
                    >
                      ×
                    </button>
                  </div>

                  {/* Navigation des catégories */}
                  <div className="cm-emoji-categories">
                    {Object.keys(emojiCategories).map(category => (
                      <button
                        key={category}
                        type="button"
                        className={`cm-emoji-category ${activeCategory === category ? 'active' : ''}`}
                        onClick={() => setActiveCategory(category)}
                      >
                        {categoryIcons[category]}
                      </button>
                    ))}
                  </div>

                  {/* Grille d'emojis */}
                  <div className="cm-emoji-grid">
                    {emojiCategories[activeCategory].map((emoji, index) => (
                      <button
                        key={index}
                        type="button"
                        className="cm-emoji"
                        onClick={() => addEmoji(emoji)}
                      >
                        {emoji}
                      </button>
                    ))}
                  </div>
                </div>
              )}

              {/* Preview du fichier sélectionné */}
              {selectedFile && (
                <div className="cm-file-preview">
                  {selectedFile.type.startsWith('image/') ? (
                    <img
                      src={URL.createObjectURL(selectedFile)}
                      alt="Preview"
                    />
                  ) : selectedFile.type.startsWith('video/') ? (
                    <video>
                      <source src={URL.createObjectURL(selectedFile)} type={selectedFile.type} />
                    </video>
                  ) : null}

                  <div className="cm-file-info">
                    <div className="cm-file-name">{selectedFile.name}</div>
                    <div className="cm-file-size">
                      {(selectedFile.size / 1024 / 1024).toFixed(2)} MB
                    </div>
                  </div>

                  <button
                    type="button"
                    className="cm-file-remove"
                    onClick={handleRemoveFile}
                  >
                    ×
                  </button>
                </div>
              )}
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}