import { Button, message, Popover } from "antd";
import { useEffect, useRef, useState } from "react";
import { io, Socket } from "socket.io-client";
import './index.scss';
import { chatHistoryList, chatroomList } from "../../api";
import type { UserInfo } from "../UpdateInfo";
import TextArea from "antd/es/input/TextArea";
import { getUserInfo } from "../../utils";
import { useLocation } from "react-router-dom";
import EmojiPicker from "@emoji-mart/react";
import data from "@emoji-mart/data";
import { UploadModal } from "./UploadModal";

interface JoinRoomPayload {
  chatroomId: number
  userId: number
}

interface SendMessagePayload {
  sendUserId: number;
  chatroomId: number;
  message: Message
}

type MessageType = 'text' | 'image' | 'file';

interface Message {
  type: MessageType
  content: string
}

type Reply = {
  type: 'sendMessage'
  userId: number
  message: ChatHistory
} | {
  type: 'joinRoom'
  userId: number
}


interface Chatroom {
  id: number;
  name: string;
  createTime: Date;
}

interface ChatHistory {
  id: number
  content: string
  type: number
  chatroomId: number
  senderId: number
  createTime: Date,
  sender: UserInfo
}

export function Chat() {
  const socketRef = useRef<Socket>(null);
  const [roomId, setChatroomId] = useState<number>();
  const [isUploadModalOpen, setUploadModalOpen] = useState(false);
  const userInfo = getUserInfo();

  useEffect(() => {
    if (!roomId) {
      return;
    }
    const socket = socketRef.current = io('http://localhost:3000');
    socket.on('connect', function () {

      const payload: JoinRoomPayload = {
        chatroomId: roomId,
        userId: userInfo.id
      }

      socket.emit('joinRoom', payload);

      socket.on('message', (reply: Reply) => {
        console.log(reply, 'reply');
        if (reply.type === 'sendMessage') {
          setChatHistory((chatHistory) => {
            return chatHistory ? [...chatHistory, reply.message] : [reply.message]
          });
          setTimeout(() => {
            document.getElementById('bottom-bar')?.scrollIntoView({ block: 'end' });
          }, 300);
        }
      });

    });
    return () => {
      socket.disconnect();
    }
  }, [roomId]);

  function sendMessage(value: string, type: MessageType = 'text') {
    if (!value) {
      return;
    }
    if (!roomId) {
      return;
    }

    const payload: SendMessagePayload = {
      sendUserId: userInfo.id,
      chatroomId: roomId,
      message: {
        type,
        content: value
      }
    }

    socketRef.current?.emit('sendMessage', payload);
  }

  const [roomList, setRoomList] = useState<Array<Chatroom>>();

  async function queryChatroomList() {
    try {
      const res = await chatroomList('');

      setRoomList(res.data.map((item: Chatroom) => {
        return {
          ...item,
          key: item.id
        }
      }));
    } catch (e: unknown) {
      if (e instanceof Error) {
        message.error(e.message);
      } else {
        message.error('系统繁忙，请稍后再试');
      }
    }
  }

  useEffect(() => {
    queryChatroomList();
  }, []);

  // 每次聊天室发生变化时，滚动到底部
  useEffect(() => {
    setTimeout(() => {
      document.getElementById('bottom-bar')?.scrollIntoView({ block: 'end' });
    }, 300);
  }, [roomId])

  const [chatHistory, setChatHistory] = useState<Array<ChatHistory>>();

  async function queryChatHistoryList(chatroomId: number) {
    try {
      const res = await chatHistoryList(chatroomId);

      setChatHistory(res.data.map((item: ChatHistory) => {
        return {
          ...item,
          key: item.id
        }
      }));
    } catch (e: unknown) {
      if (e instanceof Error) {
        message.error(e.message);
      } else {
        message.error('系统繁忙，请稍后再试');
      }
    }
  }
  const [inputText, setInputText] = useState('');

  const location = useLocation();

  useEffect(() => {
    if (location.state?.chatroomId) {
      setChatroomId(location.state?.chatroomId);
      queryChatHistoryList(location.state?.chatroomId);
    }
  }, [location.state?.chatroomId]);

  const [uploadType, setUploadType] = useState<'image' | 'file'>('image');

  return <div id="chat-container">
    <div className="chat-room-list">
      {
        roomList?.map(item => {
          return <div className={`chat-room-item ${item.id === roomId ? 'selected' : ''}`} key={item.id} data-id={item.id} onClick={() => {
            queryChatHistoryList(item.id);
            setChatroomId(item.id);
          }}>{item.name}</div>
        })
      }
    </div>
    <div className="message-list">
      {chatHistory?.map(item => {
        return <div className={`message-item ${item.senderId === userInfo.id ? 'from-me' : ''}`} data-id={item.id}>
          <div className="message-sender">
            <img src={item.sender?.headPic ?? ''} />
            <span className="sender-nickname">{item.sender?.nickName ?? ''}</span>
          </div>
          <div className="message-content">
            {item.type === 0 ? item.content : 
             item.type === 1 ? <img src={item.content} width="100" height="100" /> : 
             <div><a href={item.content} download>{item.content}</a></div>}
          </div>
        </div>
      })}
      <div id="bottom-bar" key='bottom-bar'></div>
    </div>
    <div className="message-input">
      <div className="message-type">
        <div className="message-type-item" key={1}>
          <Popover content={<EmojiPicker data={data} onEmojiSelect={(emoji: any) => {
            setInputText((inputText) => inputText + emoji.native)
          }} />} title="表情" trigger="click">
            表情
          </Popover>
        </div>
        <div className="message-type-item" key={2}>
          <div onClick={() => {
            setUploadType('image');
            setUploadModalOpen(true);
          }}>图片</div>
        </div>
        <div className="message-type-item" key={3}>
          <div onClick={() => {
            setUploadType('file');
            setUploadModalOpen(true);
          }}>文件</div>
        </div>
      </div>
      <div className="message-input-area">
        <TextArea className="message-input-box" value={inputText} onChange={(e) => {
          setInputText(e.target.value)
        }} />
        <Button className="message-send-btn" type="primary" onClick={() => {
          sendMessage(inputText)
          setInputText('');
        }}>发送</Button>
      </div>
    </div>
    <UploadModal isOpen={isUploadModalOpen} type={uploadType} handleClose={(imageSrc?: string) => {
      setUploadModalOpen(false);

      if (imageSrc) {
        sendMessage(imageSrc, uploadType);
      }
    }} />
  </div>
}
