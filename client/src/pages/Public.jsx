import React, { useEffect, useState } from 'react'
import { useSelector } from 'react-redux'
import Tweet from '../components/Tweet';
import api from '../api';
import ComposeForm from '../components/ComposeForm';

const Public = () => {
    const [tweets, setTweets] = useState([])
    const auth = useSelector(state => state.auth)

    const fetchTweets = () => {
        api.get("/tweets")
            .then(response => {
                setTweets(response.data)
            })
    }

    useEffect(() => {
        fetchTweets()
    },[])

    return (
        <div>
            <h2>Public timeline</h2>

            { auth.isLoggedIn && <ComposeForm callback={fetchTweets}/> }

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

export default Public