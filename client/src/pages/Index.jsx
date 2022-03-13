import React from 'react'
import axios from "axios"
import { useEffect, useState } from "react"
import Tweet from '../components/Tweet';

axios.defaults.baseURL = "http://localhost:8080"

const Index = () => {
    const [tweets, setTweets] = useState([])

    useEffect(() => {
        axios.get("tweets")
            .then(response => {
                setTweets(response.data)
            })
    },[])

    return (
        <div>
            <h2>Public timeline</h2>
            <ul className="messages">
                { tweets.map(tweet => {
                    return (
                        <Tweet key={tweet.Message_id} tweet={tweet}/>
                    )
                })}

                {tweets.length === 0 &&
                    <li><em>There's no message so far.</em></li>
                }
            </ul>
        </div>
    )
}

export default Index