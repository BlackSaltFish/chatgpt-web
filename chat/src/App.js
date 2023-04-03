import './App.css';
import Chat, {Bubble, Progress, toast, useMessages} from '@chatui/core';
import '@chatui/core/dist/index.css';
import '@chatui/core/es/styles/index.less';
import React, {useState} from 'react';
import './chatui-theme.css';
import axios from "axios";
import ReactMarkdown from 'react-markdown'

import CopyToClipboard from 'react-copy-to-clipboard';
import { Button } from '@chatui/core';

import {Prism as SyntaxHighlighter} from 'react-syntax-highlighter'
import {darcula} from 'react-syntax-highlighter/dist/esm/styles/prism'

const defaultQuickReplies = [
    {
        name: '清空会话',
        isNew: true,
        isHighlight: true,
    }
];


const initialMessages = [
    {
        type: 'text',
        content: {text: '您好,请问有什么可以帮您'},
        user: {avatar: '//gitclone.com/download1/gitclone.png'},
    },
];

let userAvatar = '//gitclone.com/download1/gitclone.png';


let chatContext = [];

axios.post('avatar',
   ).then((response) => {
    userAvatar = response.data.avatar
    }
).catch(err => {
    // 错误处理
    toast.fail("请求出错，" + err.response.data.errorMsg)
});


function App() {
    const {messages, appendMsg, setTyping} = useMessages(initialMessages);
    const [percentage, setPercentage] = useState(0);


    // clearQuestion 清空文本特殊字符
    function clearQuestion(requestText) {
        requestText = requestText.replace(/\s/g, "");
        const punctuation = ",.;!?，。！？、…";
        const runeRequestText = requestText.split("");
        const lastChar = runeRequestText[runeRequestText.length - 1];
        if (punctuation.indexOf(lastChar) < 0) {
            requestText = requestText + "。";
        }
        return requestText
    }

    // clearQuestion 清空文本换行符号
    function clearReply(reply) {
        // TODO 清洗回复特殊字符
        return reply
    }

    function handleSend(type, val) {
        if (percentage > 0) {
            toast.fail('正在等待上一次回复，请稍后')
            return;
        }
        if (type === 'text' && val.trim()) {
            appendMsg({
                type: 'text',
                content: {text: val},
                position: 'left',
                user: {avatar: userAvatar},
            });

            setTyping(true);
            setPercentage(10);
            onGenCode(val);
        }
    }

    function renderMessageContent(msg) {
        const {type, content} = msg;

        switch (type) {
            case 'text':
                let text = content.text
                return <Bubble>
                    <ReactMarkdown children={text} components={{
                        code({node, inline, className, children, ...props}) {
                            const match = /language-(\w+)/.exec(className || '')
                            return !inline && match ? (
                                <div>
                                <CopyToClipboard text={children}>
                                    <Button>复制代码</Button>
                                </CopyToClipboard>

                                <SyntaxHighlighter
                                    children={String(children).replace(/\n$/, '')}
                                    style={darcula}
                                    language={match[1]}
                                    PreTag="div"
                                    showLineNumbers
                                    {...props}
                                />
                                </div>
                            ) : (
                                <code className={className} {...props}>
                                    {children}
                                </code>
                            )
                        }
                    }}  />
                </Bubble>;
            default:
                return null;
        }
    }

    function handleQuickReplyClick(item) {
        if (item.name === "清空会话") {
            window.location.reload()
        }
    }

    function onGenCode(question) {
        question = clearQuestion(question)
        chatContext.push({
            "role": "user",
            "content": question,
        })


        let url = "completion"

        axios.post(url,
            {
                "messages": chatContext,
            }).then((response) => {
                let reply = clearReply(response.data.data.reply)
                appendMsg({
                    type: 'text',
                    content: {text: reply},
                    user: {avatar: '//gitclone.com/download1/gitclone.png'},
                });
                chatContext = response.data.data.messages
                console.log(chatContext)
                setPercentage(0);
            }
        ).catch(err => {
            // 错误处理
            toast.fail("请求出错，" + err.response.data.errorMsg)
        });
    }

    return (
        <div style={{height: 'calc(100vh - 10px)', marginTop: '-5px'}}>
            <Chat
                navbar={{
                    leftContent: {
                        icon: 'chevron-left',
                        title: 'Back',
                    },
                    rightContent: [
                        {
                            icon: 'apps',
                            title: 'Applications',
                        },
                        {
                            icon: 'ellipsis-h',
                            title: 'More',
                        },
                    ],
                    title: '基于ChatGPT的AI助手',
                }}
                messages={messages}
                renderMessageContent={renderMessageContent}
                quickReplies={defaultQuickReplies}
                onQuickReplyClick={handleQuickReplyClick}
                onSend={handleSend}
            />
            <Progress value={percentage}/>
        </div>
    );
}

export default App;
