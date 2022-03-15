import React from 'react'
import Tweet from "./Tweet"

const Tweets = ({ tweets }) => {
    return (
        <ul className="space-y-3">
            { tweets.map(tweet => {
                return (
                    <Tweet key={tweet.Message_id} tweet={tweet}/>
                )
            })}

            {tweets.length === 0 &&
                <li><em>There's no message so far.</em></li>
            }
        </ul>
    )
}

export default Tweets