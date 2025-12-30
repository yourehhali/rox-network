"use client";

import Link from 'next/link';
import { useEffect } from 'react';
import styles from './error.module.css';

export default function Error({ error, reset }) {
  useEffect(() => {
    console.error("Application Error:", error);
  }, [error]);

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Oops!</h1>
      <h2 className={styles.subtitle}>Something went wrong</h2>
      <p className={styles.text}>
        We encountered an unexpected error. Please try again or return to the home page.
      </p>
      <div className={styles.actions}>
        <button onClick={() => reset()} className={styles.button}>
          Try Again
        </button>
        <Link href="/" className={styles.outlineButton}>
          Back to Home
        </Link>
      </div>
    </div>
  );
}
