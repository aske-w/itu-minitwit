import React, { useState } from 'react'
import api from "../api"
import { login } from '../reducers/auth'
import { useNavigate } from "react-router-dom";
import { useDispatch } from 'react-redux'

const Signin = () => {
	const [form, setForm] = useState({
        username: "",
        password: "",
    })

	const navigate = useNavigate()
	const dispatch = useDispatch()

    const handleChange = (event) => {
        setForm({...form, [event.target.name]: event.target.value})
    }

    function handleSubmit(e) {
        e.preventDefault()

		api.post("signin", form).then(response => {
			dispatch(login(response.data))

			navigate("/")
        })
    }

  	return (
		<div>
			<h2>Sign in</h2>
			<form onSubmit={handleSubmit}>
				<dl>
					<dt>Username:</dt>
					<dd><input type="text" name="username" onChange={handleChange}/></dd>

					<dt>Password:</dt>
					<dd><input type="password" name="password" onChange={handleChange}/></dd>
				</dl>

				<div className="actions">
					<input type="submit" value="Sign In"/>
				</div>
			</form>
		</div>
  	)
}

export default Signin