import React, { useState } from 'react'
import api from '../api'

const Signup = () => {
    const [form, setForm] = useState({
        username: "",
        email: "",
        password: "",
    })

    const handleChange = (event) => {
        setForm({...form, [event.target.name]: event.target.value})
    }

    function handleSubmit(e) {
        e.preventDefault()

        api.post("/signup", form)
    }

    return (
        <div>
            <h2>Sign up</h2>

            <form action="POST" onSubmit={handleSubmit}>
                <dl>
                    <dt>Username:</dt>
                    <dd>
                        <input type="text" name="username" size="30" onChange={handleChange}/>
                    </dd>

                    <dt>E-Mail:</dt>
                    <dd>
                        <input type="text" name="email" size="30" onChange={handleChange}/>
                    </dd>

                    <dt>Password:</dt>
                    <dd>
                        <input type="password" name="password" size="30" onChange={handleChange}/>
                    </dd>
                </dl>

                <div className="actions">
                    <input type="submit" value="Sign Up"/>
                </div>
            </form>
        </div>
    )
}

export default Signup