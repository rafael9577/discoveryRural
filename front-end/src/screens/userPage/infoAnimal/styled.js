import { StyleSheet } from "react-native";

export default StyleSheet.create({
  safeArea: {
    backgroundColor: "#0A5133",
    flex: 1,
    alignItems: "center",
  },

  // --- Cabeçalho ---
  containerCabecalho: {
    flexDirection: "row",
    alignItems: "center",
    marginTop: 20,
    paddingHorizontal: 20,
    width: "100%",
  },
  tituloCabecalho: {
    flex: 1,
    textAlign: "center",
    fontSize: 22,
    fontWeight: "bold",
    color: "#FFFFFF",
    marginRight: 28,
  },

  // --- Cartão Principal ---
  suporteCartao: {
    width: "85%",
    height: "75%",
    alignSelf: "center",
    marginTop: 20,
  },
  containerPrincipal: {
    width: "100%",
    height: "100%",
    backgroundColor: "#D9D9D9",
    borderRadius: 16,
    borderWidth: 1.5,
    borderColor: "#000000",
    overflow: "hidden",
  },

  // --- Bloco Superior Interno ---
  blocoSuperior: {
    flexDirection: "row",
    padding: 15,
  },
  imagemAnimal: {
    width: 90,
    height: 90,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: "#A0A0A0",
  },
  divisorVertical: {
    width: 2,
    backgroundColor: "#000000",
    marginHorizontal: 12,
  },
  containerInformacoes: {
    flex: 1,
  },
  nomeAnimal: {
    fontWeight: "bold",
    textTransform: "uppercase",
    fontSize: 14,
    color: "#000000",
  },
  divisorHorizontal: {
    height: 2,
    backgroundColor: "#000000",
    marginVertical: 4,
  },
  nomeCientifico: {
    fontStyle: "italic",
    fontSize: 12,
    marginBottom: 6,
    color: "#222222",
  },
  textoInformacao: {
    fontSize: 12,
    fontWeight: "600",
    color: "#000000",
    marginTop: 2,
  },
});