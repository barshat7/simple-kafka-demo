import React from 'react'

const TextInput = ({onChange}) => {
    return (
        <div>
            <input type = 'text' placeholder='Your Event Name...' onChange = {(e) => onChange(e)} />
        </div>
    )
}

export default TextInput
