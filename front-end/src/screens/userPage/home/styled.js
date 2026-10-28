import { StyleSheet } from "react-native";

export default StyleSheet.create({
  safeArea: {
    backgroundColor: "#0A5133",
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    paddingVertical: 20,
  },
  cardContainer: {
    width: "80%",
    height: "60%", // Reduzido de 75% para caber o card inferior com folga
    borderWidth: 5,
    borderColor: "#2E7D32",
    borderRadius: 14,
    overflow: "hidden",
    alignSelf: "center",
    marginBottom: 20, // Garante o espaçamento direto abaixo do primeiro card
  },
  cardContainer2: {
    width: "80%",
    height: "15%", // Ajustado para proporcionalidade
    borderWidth: 5,
    borderColor: "#faf9f9",
    borderRadius: 14,
    overflow: "hidden",
    alignSelf: "center",
  },
  cardBackground: {
    width: "100%",
    height: "100%",
    justifyContent: "flex-end",
  },
  imageStyle: {
    borderRadius: 8,
  },
  overlay: {
    backgroundColor: "rgba(0, 0, 0, 0.4)",
    padding: 16,
    width: "100%",
  },
  overlay2: {
    backgroundColor: "rgba(248, 4, 4, 0.88)",
    width: "100%",
    height: "100%",
    justifyContent: "center",
    alignItems: "center",
  },
  cardText: {
    color: "#FFFFFF",
    fontSize: 22,
    fontWeight: "bold",
  },
});