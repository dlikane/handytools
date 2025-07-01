import librosa
import json
import sys

def detect_beats(audio_path, output_json='beats.json'):
    y, sr = librosa.load(audio_path)
    tempo, beat_frames = librosa.beat.beat_track(y=y, sr=sr)
    tempo = float(tempo)  # ✅ Ensure it's a native float
    beat_times = librosa.frames_to_time(beat_frames, sr=sr)

    with open(output_json, 'w') as f:
        json.dump(beat_times.tolist(), f, indent=2)

    print(f"Detected {len(beat_times)} beats at ~{tempo:.2f} BPM.")
    print(f"Saved to {output_json}")

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Usage: python detect_beats.py audio_file.wav")
        sys.exit(1)

    audio_file = sys.argv[1]
    output = sys.argv[2] if len(sys.argv) > 2 else 'beats.json'
    detect_beats(audio_file, output)
