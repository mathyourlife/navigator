import React, { useState, useEffect } from 'react';
import Modal from 'react-modal'; // Import your preferred modal library
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Button } from '@mui/material';

const Skill = () => {
  const [skills, setSkills] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [newSkill, setNewSkill] = useState({ name: '', description: ''});

  useEffect(() => {
    fetch('/api/skill')
      .then(res => res.json())
      .then(data => setSkills(data.skills));
  }, []);

  const handleInputChange = (event) => {
    setNewSkill({ ...newSkill, [event.target.name]: event.target.value });
  };

  const handleCreateSkill = async () => {
      fetch('/api/skill', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({skill: newSkill})
      })
        .then(res => res.json())
        .then(data => {
            setSkills([...skills, data.skill]);
            setShowModal(false);
                    setNewSkill({ name: '', description: ''});
        })
        .catch(error => console.error('Error creating skill:', error));
  };

  const handleDeleteSkill = async (skillId) => {
    try {
      const response = await fetch(`/api/skill/${skillId}`, {
        method: 'DELETE',
      });

      if (response.ok) {
        // Refresh the skill list after successful deletion
        fetch('/api/skill')
          .then(res => res.json())
          .then(data => setSkills(data.skills));
      } else {
        console.error('Failed to delete skill:', response.status);
      }
    } catch (error) {
      console.error('Error deleting skill:', error);
    }
  };

  const handleEditSkill = async (skillId, updatedSkill) => {
    try {
      const response = await fetch(`/api/skill/${skillId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updatedSkill),
      });

      if (response.ok) {
        // Refresh the skill list after successful update
        fetch('/api/skill')
          .then(res => res.json())
          .then(data => setSkills(data.skills));
      } else {
        console.error('Failed to update skill:', response.status);
      }
    } catch (error) {
      console.error('Error updating skill:', error);
    }
  };

  return (
    <div className="container">
      <h1>Skill List</h1>

      {/* Modal for creating new skills */}
      <Modal isOpen={showModal} onRequestClose={() => setShowModal(false)}>
        <h2>Create New Skill</h2>
        <input
          type="text"
          name="name"
          placeholder="Skill Name"
          value={newSkill.name}
          onChange={handleInputChange}
        />
        <input
          type="text"
          name="description"
          placeholder="Description"
          value={newSkill.description}
          onChange={handleInputChange}
        />
        <button onClick={handleCreateSkill}>Create</button>
        <button onClick={() => setShowModal(false)}>Cancel</button>
      </Modal>

      {/* Button to open the modal */}
      <button onClick={() => setShowModal(true)}>Create New Skill</button>

      {/* Table to display skills */}
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Description</TableCell>
              <TableCell>Action</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {skills.length === 0 && (
              <TableRow>
                <TableCell colSpan={3} align="center">No skills found</TableCell>
                </TableRow>
            )}
            {skills.map((skill) => (
              <TableRow key={skill.skill_id}>
                <TableCell>{skill.name}</TableCell>
                <TableCell>{skill.description}</TableCell>
                <TableCell>
                  <button onClick={() => handleDeleteSkill(skill.skill_id)}>Delete</button>
                  <button onClick={() => handleEditSkill(skill.skill_id, { name: skill.name })}>Edit</button> {/* Placeholder, you'll need to implement the edit functionality */}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
};

export default Skill;
