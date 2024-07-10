import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';

const MovieList = () => {
  const [movies, setMovies] = useState([]);

  useEffect(() => {
    fetchMovies();
  }, []);

  const fetchMovies = () => {
    axios.get('http://localhost:8000/movies')
      .then(response => {
        setMovies(response.data || []);
      })
      .catch(error => {
        console.error('There was an error fetching the movies!', error);
        setMovies([]);  // Ensure movies is an array even if the fetch fails
      });
  };

  const deleteMovie = (id) => {
    axios.delete(`http://localhost:8000/movies/${id}`)
      .then(() => {
        fetchMovies();
      })
      .catch(error => {
        console.error('There was an error deleting the movie!', error);
      });
  };

  return (
    <div>
      <h2>Movie List</h2>
      {movies.length === 0 ? (
        <div>No movies available</div>
      ) : (
        <ul>
          {movies.map(movie => (
            <li key={movie.id}>
              <Link to={`/movies/${movie.id}`}>
                {movie.title} - {movie.director.first_name} {movie.director.last_name}
              </Link>
              <button onClick={() => deleteMovie(movie.id)}>Delete</button>
              <Link to={`/edit-movie/${movie.id}`}>Edit</Link>
            </li>
          ))}
        </ul>
      )}
      <Link to="/add-movie">Add a new movie</Link>
    </div>
  );
};

export default MovieList;
