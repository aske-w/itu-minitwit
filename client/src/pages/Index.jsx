import React, { useEffect, useState } from 'react'
import Tweet from '../components/Tweet';
import { useSelector } from 'react-redux'
import { useNavigate } from "react-router-dom";
import api from '../api';
import ComposeForm from '../components/ComposeForm';

const Index = () => {
    const [tweets, setTweets] = useState([])
    const auth = useSelector(state => state.auth)
    const navigate = useNavigate()

    const fetchTweets = () => {
        api.get("/timeline")
            .then(response => {
                setTweets(response.data)
            })
    }

    useEffect(() => {
        if ( ! auth.isLoggedIn) {
            navigate("/public")

            return
        }

        fetchTweets()
    }, [])

    return (
        <div>
            <h2>My timeline</h2>

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

export default Index