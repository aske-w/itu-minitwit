import React, { useEffect, useState } from 'react'
import Tweets from '../components/Tweets';
import { useSelector } from 'react-redux'
import { useNavigate, useParams } from "react-router-dom";
import api from '../api';

const Profile = (props) => {
    const [user, setUser] = useState(null)
    const [tweets, setTweets] = useState([])
    const auth = useSelector(state => state.auth)
    const navigate = useNavigate()
    const { username } = useParams()

    const fetchUser = () => {
        api.get(`/users/${username}`)
            .then(response => {
                setUser(response.data)
            })
            .catch(error => {
                navigate("/public")
            })
    }

    const fetchTweets = () => {
        api.get(`/users/${user.username.toLowerCase()}/tweets`)
            .then(response => {
                setTweets(response.data)
            })
    }

    useEffect(() => {
        fetchUser()
    }, [])

    useEffect(() => {
        if ( user) {
            fetchTweets()
        }
    }, [user])



    if ( ! user) {
        return null
    }

    return (
        <div>
            <h2 className="text-2xl font-semibold mb-4">{ user.username }'s timeline</h2>

            <Tweets tweets={tweets}/>
        </div>
    )
}

export default Profile