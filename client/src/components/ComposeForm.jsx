import React, { useState, useRef } from 'react'
import { useSelector} from 'react-redux'
import api from "../api"

const ComposeForm = ({ callback }) => {
	const textInput = useRef(null)
	const [form, setForm] = useState({
		text: "",
	})
	const auth = useSelector(state => state.auth)

	const handleChange = (event) => {
        setForm({...form, [event.target.name]: event.target.value})
    }

	function handleSubmit(e) {
        e.preventDefault()

		api.post("tweets", form).then(() => {
			setForm({text: ""})
			textInput.current.value = ""

			callback()
		})
    }

	return (
		<div className="twitbox">
			<h3>What's on your mind {auth.user.username}?</h3>
			<form onSubmit={handleSubmit}>
				<p>
					<input ref={textInput} type="text" name="text" size="60" onChange={handleChange}/>
					<input type="submit" value="Share" />
				</p>
			</form>
		</div>
	)
}

export default ComposeForm