import React, { useEffect, useState } from 'react'
import Tweets from '../components/Tweets';
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
            <h2 className="text-2xl font-semibold mb-4">My timeline</h2>

            { auth.isLoggedIn && <ComposeForm callback={fetchTweets}/> }

            <Tweets tweets={tweets}/>
        </div>
    )
}

export default Index