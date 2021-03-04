import React from 'react'
import Button from './Button';

const VoteButtons = ({type, candidates, onVote}) => {
    return (
        <>
            {
                candidates.map(candidate => {
                    return <Button key = {candidate.uniqueID} text = {`${type} ${candidate.uniqueID}`} color = {candidate.color} onClick = {() => onVote(candidate.uniqueID)} />
                })
            }
        </>
    )
}

export default VoteButtons
